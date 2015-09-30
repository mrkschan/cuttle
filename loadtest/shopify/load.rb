#!/usr/bin/ruby

require 'logger'
require 'net/http'

def generate_load(ins, thread_count)
  logger = Logger.new(STDERR)

  queue = Queue.new
  ins.map { |i| queue << i }

  threads = thread_count.times.map do
    Thread.new do
      while !queue.empty? && i = queue.pop
        res = `bash -c 'source env.sh && env/bin/python apicall.py'`
        if res.include? 'Too Many Requests'
          logger.warn "#{res}"
        end
      end
    end
  end

  threads.each(&:join)
end


if __FILE__ == $0
  generate_load([1] * 160, 32)
end

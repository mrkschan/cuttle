#!/usr/bin/ruby

require 'logger'
require 'net/http'

def generate_load(urls, thread_count)
  logger = Logger.new(STDERR)

  queue = Queue.new
  urls.map { |url| queue << url }

  threads = thread_count.times.map do
    Thread.new do
      while !queue.empty? && url = queue.pop
        res = `curl -k -s "#{url}"`
        if res.include? 'OVER_QUERY_LIMIT'
          logger.warn "#{res}"
        end
      end
    end
  end

  threads.each(&:join)
end


if __FILE__ == $0
  api_key = ENV['API_KEY']
  urls = ["https://maps.googleapis.com/maps/api/geocode/json?address=1600+Amphitheatre+Parkway,+Mountain+View,+CA&key=#{api_key}",
          "https://maps.googleapis.com/maps/api/geocode/json?address=White+House&key=#{api_key}",
          "https://maps.googleapis.com/maps/api/geocode/json?address=Google&key=#{api_key}",
          "https://maps.googleapis.com/maps/api/geocode/json?address=Apple&key=#{api_key}",
          ]
  generate_load(urls * 32, 32)
end

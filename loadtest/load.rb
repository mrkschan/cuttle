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
        # TODO(mrkschan): Should not use `curl -k`
        # Ref - https://github.com/elazarl/goproxy/issues/73
        res = `curl -k -s "#{url}"`
        if res.include? 'OVER_QUERY_LIMIT'
          logger.warn "#{url} - #{res}"
        end
      end
    end
  end

  threads.each(&:join)
end


if __FILE__ == $0
  googleapis_key = 'AIzaSyCQ1dg6d95v88xorMRFUrLw1vGi8j8_l3Y'
  urls = ["https://maps.googleapis.com/maps/api/geocode/json?address=1600+Amphitheatre+Parkway,+Mountain+View,+CA&key=#{googleapis_key}",
          "https://maps.googleapis.com/maps/api/geocode/json?address=White+House&key=#{googleapis_key}",
          "https://maps.googleapis.com/maps/api/geocode/json?address=Google&key=#{googleapis_key}",
          "https://maps.googleapis.com/maps/api/geocode/json?address=Apple&key=#{googleapis_key}",
          ]
  generate_load(urls * 32, 32)
end

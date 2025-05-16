local wrk = require("wrk")

math.randomseed(os.time())

function generate_random_title()
  local length = math.random(5, 15)
  local chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
  local title = ""
  
  for i = 1, length do
    local rand = math.random(1, #chars)
    title = title .. string.sub(chars, rand, rand)
  end
  
  return "Playlist_" .. title
end

local count = 1

function response()
    if count == 1000 then
        wrk.thread:stop()
        count = 1
    end
    count = count + 1
end

local boundary = "----WebKitFormBoundary" .. string.sub(tostring(math.random()), 3, 14)

local title = generate_random_title()
local body = "--" .. boundary .. "\r\n" ..
             "Content-Disposition: form-data; name=\"title\"\r\n\r\n" ..
             title .. "\r\n" ..
             "--" .. boundary .. "--\r\n"

wrk.method = "POST"
wrk.body = body
wrk.headers["Content-Type"] = "multipart/form-data; boundary=" .. boundary
wrk.headers["Content-Length"] = string.len(wrk.body)


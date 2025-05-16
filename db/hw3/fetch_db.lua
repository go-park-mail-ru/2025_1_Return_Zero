local wrk = require("wrk")

wrk.method = "GET"
wrk.headers["Content-Type"] = "application/json"
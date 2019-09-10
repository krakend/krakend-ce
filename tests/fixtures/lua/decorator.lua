function post_proxy_decorator( resp )
	local responseData = resp:data()
	local responseContent = responseData:get("source_result")
	local message = responseContent:get("message")

	local c = string.match(message, "Successfully")

	if not not c
	then
		responseData:set("result", "success")
	else
		responseData:set("result", "failed")
	end
end
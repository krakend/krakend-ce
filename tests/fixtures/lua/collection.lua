function pre_backend( req )
	local authHeader = req:headers("Authorization")
	local first = string.find(authHeader, "%.")
	local last = string.find(authHeader, "%.", first+1)
	local rawData = string.sub(authHeader, first+1, last-1)
	local decoded = from_base64(rawData)
	local jwtData = json_parse(decoded)

	req:url(req:url() .. jwtData["id"])
end

function post_proxy( resp )
	local responseData = resp:data()
	local data = {}
	local col = responseData:get("collection")

	local size = col:len()
	responseData:set("total", size)

	local paths = {}
	for i=0,size-1 do
		local element = col:get(i)
		local t = element:get("path")
		table.insert(paths, t)
	end
	responseData:set("paths", paths)
	responseData:set("collection", {})
end

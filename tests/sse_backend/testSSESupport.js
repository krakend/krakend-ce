/**
  How to Use This Server:

  Regular HTTP endpoint:
  Access http://localhost:5555/api/data to get a simple JSON response

  SSE endpoint that closes automatically:
  Access http://localhost:5555/api/events to connect to the SSE stream
  Connection will close after 5 seconds (default)
  You can customize the timeout with http://localhost:5555/api/events?closeAfter=10000 (in milliseconds)


  Health check:
  Access http://localhost:5555/health to verify the server is running
*/

const http = require('http');
const url = require('url');

const PORT = 5555;

// Create HTTP server
const server = http.createServer((req, res) => {
  // Parse URL and query parameters
  const parsedUrl = url.parse(req.url, true);
  const pathname = parsedUrl.pathname;
  const query = parsedUrl.query;
  
  console.log(`${req.method} request received for ${pathname}`);
  
  // Regular HTTP GET endpoint
  if (req.method === 'GET' && pathname === '/api/data') {
    res.writeHead(200, {
      'Content-Type': 'application/json',
      'Access-Control-Allow-Origin': '*'
    });
    
    const responseData = {
      message: 'This is a regular HTTP response',
      timestamp: new Date().toISOString()
    };
    
    res.end(JSON.stringify(responseData));
  }
  
  // SSE endpoint
  else if (req.method === 'GET' && pathname === '/api/events') {
    // Set headers for SSE
    res.writeHead(200, {
      'Content-Type': 'text/event-stream',
      'Cache-Control': 'no-cache',
      'Connection': 'keep-alive',
      'Access-Control-Allow-Origin': '*'
    });
    
    // Log connection
    const clientIp = req.socket.remoteAddress;
    console.log(`SSE Client connected from ${clientIp}`);
    
    // Get connection timeout from query param or use default (5 seconds)
    const closeAfter = parseInt(query.closeAfter) || 5000;
    console.log(`Will close connection after ${closeAfter}ms`);
    
    // Send initial connection event
    const initialEvent = {
      type: 'connected',
      message: 'Connection established',
      closeAfter: closeAfter
    };
    res.write(`data: ${JSON.stringify(initialEvent)}\n\n`);
    
    // Send events every second
    let count = 0;
    const intervalId = setInterval(() => {
      const eventData = {
        id: count,
        type: 'update',
        message: `Event update ${count}`,
        timestamp: new Date().toISOString()
      };
      
      res.write(`id: ${count}\n`);
      res.write(`data: ${JSON.stringify(eventData)}\n\n`);
      
      count++;
    }, 1000);
    
    // Handle client disconnection
    req.on('close', () => {
      console.log('Client disconnected (connection closed by client)');
      clearInterval(intervalId);
    });
    
    // Close connection after specified time
    setTimeout(() => {
      console.log('Server is intentionally closing the SSE connection');
      clearInterval(intervalId);
      
      // Send final event before closing
      const finalEvent = {
        type: 'closing',
        message: 'Server is closing the connection',
        timestamp: new Date().toISOString()
      };
      
      res.write(`event: closing\n`);
      res.write(`data: ${JSON.stringify(finalEvent)}\n\n`);
      
      // End the response which closes the connection
      res.end();
    }, closeAfter);
  }
  
  // Health check endpoint
  else if (req.method === 'GET' && pathname === '/health') {
    res.writeHead(200, { 'Content-Type': 'text/plain' });
    res.end('Server is running');
  }
  
  // Handle 404 Not Found
  else {
    res.writeHead(404, { 'Content-Type': 'text/plain' });
    res.end('Not Found');
  }
});

// Start the server
server.listen(PORT, () => {
  console.log(`Test server running on port ${PORT}`);
  console.log(`- Regular HTTP endpoint: http://localhost:${PORT}/api/data`);
  console.log(`- SSE endpoint: http://localhost:${PORT}/api/events`);
  console.log(`- SSE with custom timeout: http://localhost:${PORT}/api/events?closeAfter=10000`);
});
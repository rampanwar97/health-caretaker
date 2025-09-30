let ws;
let endpoints = new Map();

// Initialize WebSocket connection
function initWebSocket() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    ws = new WebSocket(protocol + '//' + window.location.host + '/ws');
    
    ws.onopen = function() {
        console.log('WebSocket connected');
    };
    
    ws.onmessage = function(event) {
        const endpoint = JSON.parse(event.data);
        endpoints.set(endpoint.id, endpoint);
        renderEndpoints();
    };
    
    ws.onclose = function() {
        console.log('WebSocket disconnected, reconnecting...');
        setTimeout(initWebSocket, 3000);
    };
    
    ws.onerror = function(error) {
        console.error('WebSocket error:', error);
    };
}

// Load initial endpoints
async function loadEndpoints() {
    try {
        const response = await fetch('/api/endpoints');
        const data = await response.json();
        endpoints.clear();
        data.forEach(endpoint => {
            endpoints.set(endpoint.id, endpoint);
        });
        renderEndpoints();
    } catch (error) {
        console.error('Error loading endpoints:', error);
    }
}

// Render endpoints to the UI
function renderEndpoints() {
    const container = document.getElementById('endpointsContainer');
    
    if (endpoints.size === 0) {
        container.innerHTML = '<div class="loading">No endpoints configured. Add one above to get started!</div>';
        return;
    }
    
    const endpointsArray = Array.from(endpoints.values());
    let html = '<div class="endpoints-grid">';
    endpointsArray.forEach(endpoint => {
        html += '<div class="endpoint-card ' + endpoint.status + '">';
        html += '<div class="endpoint-header">';
        html += '<div class="endpoint-name">' + endpoint.name + '</div>';
        html += '<div class="status-badge status-' + endpoint.status + '">' + endpoint.status + '</div>';
        html += '</div>';
        html += '<div class="endpoint-url">' + endpoint.method + ' ' + endpoint.url + '</div>';
        html += '<div class="endpoint-details">';
        html += '<div class="detail-item"><div class="detail-label">Status Code</div><div class="detail-value">' + (endpoint.statusCode || 'N/A') + '</div></div>';
        html += '<div class="detail-item"><div class="detail-label">Response Time</div><div class="detail-value">' + (endpoint.responseTime || 0) + 'ms</div></div>';
        html += '<div class="detail-item"><div class="detail-label">Last Check</div><div class="detail-value">' + formatTime(endpoint.lastCheck) + '</div></div>';
        html += '<div class="detail-item"><div class="detail-label">Interval</div><div class="detail-value">' + endpoint.interval + 's</div></div>';
        html += '</div>';
        if (endpoint.error) {
            html += '<div class="error-message">' + endpoint.error + '</div>';
        }
        html += '<div class="endpoint-actions">';
        html += '<button class="btn btn-success" onclick="checkEndpoint(\'' + endpoint.id + '\')">Check Now</button>';
        html += '<button class="btn btn-danger" onclick="removeEndpoint(\'' + endpoint.id + '\')">Remove</button>';
        html += '</div>';
        html += '</div>';
    });
    html += '</div>';
    container.innerHTML = html;
}

// Format time for display
function formatTime(timeString) {
    const date = new Date(timeString);
    const now = new Date();
    const diff = now - date;
    
    if (diff < 60000) { // Less than 1 minute
        return 'Just now';
    } else if (diff < 3600000) { // Less than 1 hour
        return Math.floor(diff / 60000) + 'm ago';
    } else if (diff < 86400000) { // Less than 1 day
        return Math.floor(diff / 3600000) + 'h ago';
    } else {
        return date.toLocaleDateString();
    }
}

// Add new endpoint
document.getElementById('endpointForm').addEventListener('submit', async function(e) {
    e.preventDefault();
    
    const formData = new FormData(e.target);
    const endpoint = {
        name: formData.get('name'),
        url: formData.get('url'),
        method: formData.get('method'),
        interval: parseInt(formData.get('interval')),
        timeout: parseInt(formData.get('timeout'))
    };
    
    try {
        const response = await fetch('/api/endpoints', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(endpoint)
        });
        
        if (response.ok) {
            e.target.reset();
            loadEndpoints(); // Reload to get the new endpoint
        } else {
            alert('Error adding endpoint');
        }
    } catch (error) {
        console.error('Error adding endpoint:', error);
        alert('Error adding endpoint');
    }
});

// Check endpoint manually
async function checkEndpoint(id) {
    try {
        await fetch('/api/endpoints/' + id + '/check', {
            method: 'POST'
        });
    } catch (error) {
        console.error('Error checking endpoint:', error);
    }
}

// Remove endpoint
async function removeEndpoint(id) {
    if (!confirm('Are you sure you want to remove this endpoint?')) {
        return;
    }
    
    try {
        const response = await fetch('/api/endpoints/' + id, {
            method: 'DELETE'
        });
        
        if (response.ok) {
            endpoints.delete(id);
            renderEndpoints();
        } else {
            alert('Error removing endpoint');
        }
    } catch (error) {
        console.error('Error removing endpoint:', error);
        alert('Error removing endpoint');
    }
}

// Initialize the application
document.addEventListener('DOMContentLoaded', function() {
    initWebSocket();
    loadEndpoints();
});

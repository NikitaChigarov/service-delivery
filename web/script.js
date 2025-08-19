function getOrder() {
    const orderId = document.getElementById('orderId').value.trim();
    if (!orderId) {
        showError('Please enter an Order ID');
        return;
    }

    // Show loading, hide error and previous results
    document.getElementById('loading').style.display = 'block';
    document.getElementById('error').style.display = 'none';
    document.getElementById('orderInfo').style.display = 'none';

    fetch(`/api/orders/${orderId}`)
        .then(response => {
            if (!response.ok) {
                return response.json().then(err => {
                    throw new Error(err.message || 'Failed to fetch order');
                });
            }
            return response.json();
        })
        .then(order => {
            displayOrder(order);
        })
        .catch(error => {
            showError(error.message);
        })
        .finally(() => {
            document.getElementById('loading').style.display = 'none';
        });
}

function showError(message) {
    const errorElement = document.getElementById('error');
    errorElement.textContent = message;
    errorElement.style.display = 'block';
}

function displayOrder(order) {
    // Display order details
    const orderDetails = document.getElementById('orderDetails');
    orderDetails.innerHTML = `
        <div class="detail-row">
            <div class="detail-label">Order UID:</div>
            <div class="detail-value">${order.order_uid}</div>
        </div>
        <div class="detail-row">
            <div class="detail-label">Track Number:</div>
            <div class="detail-value">${order.track_number}</div>
        </div>
        <div class="detail-row">
            <div class="detail-label">Entry:</div>
            <div class="detail-value">${order.entry}</div>
        </div>
        <div class="detail-row">
            <div class="detail-label">Date Created:</div>
            <div class="detail-value">${order.date_created}</div>
        </div>
        <div class="detail-row">
            <div class="detail-label">Customer ID:</div>
            <div class="detail-value">${order.customer_id}</div>
        </div>
    `;

    // Display delivery details
    const deliveryDetails = document.getElementById('deliveryDetails');
    deliveryDetails.innerHTML = `
        <div class="detail-row">
            <div class="detail-label">Name:</div>
            <div class="detail-value">${order.delivery.name}</div>
        </div>
        <div class="detail-row">
            <div class="detail-label">Phone:</div>
            <div class="detail-value">${order.delivery.phone}</div>
        </div>
        <div class="detail-row">
            <div class="detail-label">Address:</div>
            <div class="detail-value">${order.delivery.city}, ${order.delivery.address}, ${order.delivery.region}, ${order.delivery.zip}</div>
        </div>
        <div class="detail-row">
            <div class="detail-label">Email:</div>
            <div class="detail-value">${order.delivery.email}</div>
        </div>
    `;

    // Display payment details
    const paymentDetails = document.getElementById('paymentDetails');
    paymentDetails.innerHTML = `
        <div class="detail-row">
            <div class="detail-label">Transaction:</div>
            <div class="detail-value">${order.payment.transaction}</div>
        </div>
        <div class="detail-row">
            <div class="detail-label">Amount:</div>
            <div class="detail-value">${order.payment.amount / 100} ${order.payment.currency}</div>
        </div>
        <div class="detail-row">
            <div class="detail-label">Provider:</div>
            <div class="detail-value">${order.payment.provider}</div>
        </div>
        <div class="detail-row">
            <div class="detail-label">Payment Date:</div>
            <div class="detail-value">${new Date(order.payment.payment_dt * 1000).toLocaleString()}</div>
        </div>
        <div class="detail-row">
            <div class="detail-label">Bank:</div>
            <div class="detail-value">${order.payment.bank}</div>
        </div>
    `;

    // Display items
    const itemsBody = document.getElementById('itemsBody');
    itemsBody.innerHTML = order.items.map(item => `
        <tr>
            <td>${item.name}</td>
            <td>${item.price / 100}</td>
            <td>${item.total_price / item.price}</td>
            <td>${item.total_price / 100}</td>
            <td>${item.brand}</td>
            <td>${item.status}</td>
        </tr>
    `).join('');

    // Show the order info section
    document.getElementById('orderInfo').style.display = 'block';
}

// Handle Enter key in the input field
document.getElementById('orderId').addEventListener('keypress', function(e) {
    if (e.key === 'Enter') {
        getOrder();
    }
});
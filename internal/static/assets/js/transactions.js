// Global variables for pagination
let currentPage = 1;
let limit = 10;
let selectedMerchantId = null;

// Function to get cookie value
function getCookie(name) {
  const value = `; ${document.cookie}`;
  const parts = value.split(`; ${name}=`);
  if (parts.length === 2) return parts.pop().split(';').shift();
  return null;
}

// Function to fetch transactions with optional merchant filter
async function fetchTransactions(page, merchantId = null) {
  try {
    const userId = getCookie('user_id');
    if (!userId) {
      console.error('User ID not found in cookies');
      return null;
    }

    let url = `http://localhost:8080/api/transactions/merchant/${merchantId}?page=${page}&limit=${limit}`;
    if (merchantId && merchantId !== 'all') {
      url += `&merchant_id=${merchantId}`;
    }

    const response = await fetch(url, {
      headers: {
        'accept': 'application/json',
        'Authorization': `Bearer ${getCookie('auth_token')}`,
        'X-User-Id': userId
      }
    });
    return await response.json();
  } catch (error) {
    console.error('Error fetching transactions:', error);
    return null;
  }
}

// Function to format date
function formatDate(dateString) {
  return new Date(dateString).toLocaleString();
}

// Function to get status badge class
function getStatusBadgeClass(status) {
  switch(status.toLowerCase()) {
    case 'completed':
      return 'bg-gradient-success';
    case 'pending':
      return 'bg-gradient-warning';
    case 'failed':
      return 'bg-gradient-danger';
    default:
      return 'bg-gradient-secondary';
  }
}

// Function to update transaction table
function updateTable(data) {
  const tableBody = document.getElementById('transactionTableBody');
  if (!tableBody) {
    console.error('Transaction table body not found');
    return;
  }

  tableBody.innerHTML = '';

  // Check if data and transactions array exist
  if (!data || !Array.isArray(data.transactions)) {
    console.error('Invalid data format or missing transactions array');
    return;
  }

  data.transactions.forEach(tx => {
    if (!tx) return; // Skip if transaction object is null or undefined

    const row = document.createElement('tr');
    row.innerHTML = `
      <td>
        <div class="d-flex px-2 py-1">
          <div class="d-flex flex-column justify-content-center">
            <h6 class="mb-0 text-sm">${tx.transaction_id || 'N/A'}</h6>
            <p class="text-xs text-secondary mb-0">Merchant: ${tx.merchant_id || 'N/A'}</p>
          </div>
        </div>
      </td>
      <td>
        <p class="text-xs font-weight-bold mb-0">$${tx.transaction_amount || '0.00'}</p>
        <p class="text-xs text-secondary mb-0">Program: ${tx.program_id || 'N/A'}</p>
      </td>
      <td class="align-middle text-center text-sm">
        <p class="text-xs font-weight-bold mb-0">${tx.transaction_type || 'N/A'}</p>
      </td>
      <td class="align-middle text-center text-sm">
        <span class="badge badge-sm ${getStatusBadgeClass(tx.status || 'unknown')}">${tx.status || 'Unknown'}</span>
      </td>
      <td class="align-middle text-center">
        <span class="text-secondary text-xs font-weight-bold">${formatDate(tx.transaction_date) || 'N/A'}</span>
      </td>
    `;
    tableBody.appendChild(row);
  });

  	
// Response body
// {
//   "pagination": {
//     "current_page": 1,
//     "per_page": 10,
//     "total_items": 0,
//     "total_pages": 0
//   },
//   "transactions": null
// }

  // Update pagination if pagination data exists
  if (data.pagination) {
    document.getElementById('currentPage').textContent = data.pagination.current_page || 1;
    document.getElementById('totalPages').textContent = data.pagination.total_pages || 1;
    
    // Update button states
    document.getElementById('prevPage').disabled = currentPage <= 1;
    document.getElementById('nextPage').disabled = currentPage >= (data.pagination.total_pages || 1);
  }
}

// Function to load transactions
async function loadTransactions(page, merchantId = null) {
  const data = await fetchTransactions(page, merchantId);
  if (data) {
    currentPage = page;
    updateTable(data);
  }
}

// Function to fetch merchants and populate the dropdown
async function populateMerchantDropdown() {
    try {
        const response = await fetch('http://localhost:8080/api/merchants', {
            method: 'GET',
            headers: {
                'accept': 'application/json'
            }
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const merchants = await response.json();
        const dropdown = document.getElementById('merchantDropdown');
        const dropdownMenu = dropdown.nextElementSibling;

        // Clear existing items except "All Merchants"
        dropdownMenu.innerHTML = '<li><a class="dropdown-item" href="#" data-merchant-id="all">All Merchants</a></li>';

        // Add merchants to dropdown
        merchants.forEach(merchant => {
            const li = document.createElement('li');
            li.innerHTML = `<a class="dropdown-item" href="#" data-merchant-id="${merchant.id}">${merchant.merchant_name}</a>`;
            dropdownMenu.appendChild(li);
        });

        // Add click event listeners to dropdown items
        dropdownMenu.querySelectorAll('.dropdown-item').forEach(item => {
            item.addEventListener('click', function(e) {
                e.preventDefault();
                const merchantId = this.getAttribute('data-merchant-id');
                const merchantName = this.textContent;
                dropdown.textContent = merchantName;
                selectedMerchantId = merchantId;
                loadTransactions(1, merchantId);
            });
        });

    } catch (error) {
        console.error('Error fetching merchants:', error);
    }
}

// Event listeners for pagination
document.addEventListener('DOMContentLoaded', () => {
    // Initialize merchant dropdown
    populateMerchantDropdown();

    // Add pagination event listeners
    document.getElementById('prevPage').addEventListener('click', () => {
        if (currentPage > 1) {
            loadTransactions(currentPage - 1, selectedMerchantId);
        }
    });

    document.getElementById('nextPage').addEventListener('click', () => {
        loadTransactions(currentPage + 1, selectedMerchantId);
    });

    document.getElementById('pageSize').addEventListener('change', (e) => {
        limit = parseInt(e.target.value);
        loadTransactions(1, selectedMerchantId);
    });

    // Initial load of transactions
    loadTransactions(1);
});
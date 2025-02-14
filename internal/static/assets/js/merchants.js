let currentPage = 1;
let limit = 10;

function getCookie(name) {
  const value = `; ${document.cookie}`;
  const parts = value.split(`; ${name}=`);
  if (parts.length === 2) return parts.pop().split(';').shift();
  return null;
}

async function fetchMerchants() {
  try {
    const userId = getCookie('user_id');
    if (!userId) {
      console.error('User ID not found in cookies');
      return null;
    }

    const response = await fetch(`http://localhost:8080/api/merchants/user/${userId}?page=${currentPage}&limit=${limit}`, {
      headers: {
        'accept': 'application/json',
        'Authorization': `Bearer ${getCookie('auth_token')}`,
        'X-User-Id': userId
      }
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data = await response.json();
    return data;
  } catch (error) {
    console.error('Error fetching merchants:', error);
    return null;
  }
}

function formatDate(dateString) {
  return new Date(dateString).toLocaleString();
}

function updateTable(merchants) {
  const tableBody = document.getElementById('merchantsTableBody');
  if (!tableBody) {
    console.error('Merchant table body not found');
    return;
  }

  tableBody.innerHTML = '';

  if (!merchants || !merchants.merchants || !merchants.pagination) {
    console.error('Invalid merchants data format');
    return;
  }

  merchants.merchants.forEach(merchant => {
    const row = document.createElement('tr');
    row.innerHTML = `
      <td>
        <div class="d-flex px-2 py-1">
          <div class="d-flex flex-column justify-content-center">
            <h6 class="mb-0 text-sm">${merchant.merchant_name}</h6>
            <p class="text-xs text-secondary mb-0">ID: ${merchant.id}</p>
          </div>
        </div>
      </td>
      <td>
        <p class="text-xs font-weight-bold mb-0">${merchant.merchant_type}</p>
      </td>
      <td class="align-middle text-center">
        <span class="text-secondary text-xs font-weight-bold">${formatDate(merchant.created_at)}</span>
      </td>
      <td class="align-middle text-center">
        <span class="text-secondary text-xs font-weight-bold">${formatDate(merchant.updated_at)}</span>
      </td>
      <td class="align-middle text-center">
        <span class="${merchant.status.toLowerCase() === 'active' ? 'text-success' : 'text-dark'} text-xs font-weight-bold">${merchant.status}</span>
      </td>
      <td class="align-middle text-center">
        <button class="btn btn-link text-primary mb-0 me-2" onclick="editMerchant('${merchant.id}')">
          Edit
        </button>
        <button class="btn btn-link text-danger mb-0" onclick="deleteMerchant('${merchant.id}')">
          Deactivate
        </button>
      </td>
    `;
    tableBody.appendChild(row);
  });

  // Update pagination information
  document.getElementById('currentPage').textContent = merchants.pagination.current_page;
  document.getElementById('totalPages').textContent = merchants.pagination.total_pages;

  // Update pagination button states
  document.getElementById('prevPage').disabled = currentPage <= 1;
  document.getElementById('nextPage').disabled = currentPage >= merchants.pagination.total_pages;
}

async function loadMerchants() {
  const merchants = await fetchMerchants();
  if (merchants) {
    updateTable(merchants);
  }
}

async function fetchMerchantById(merchantId) {
  try {
    const userId = getCookie('user_id');
    if (!userId) {
      console.error('User ID not found in cookies');
      return null;
    }

    const response = await fetch(`http://localhost:8080/api/merchants/${merchantId}`, {
      headers: {
        'accept': 'application/json',
        'Authorization': `Bearer ${getCookie('auth_token')}`,
        'X-User-Id': userId
      }
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    return await response.json();
  } catch (error) {
    console.error('Error fetching merchant details:', error);
    return null;
  }
}

async function editMerchant(merchantId) {
  const merchant = await fetchMerchantById(merchantId);
  if (!merchant) {
    alert('Failed to fetch merchant details');
    return;
  }

  // Update modal title and submit button
  document.getElementById('addMerchantModalLabel').textContent = 'Edit Merchant';
  const submitButton = document.querySelector('#addMerchantForm button[type="submit"]');
  submitButton.textContent = 'Update Merchant';

  // Populate form fields
  document.getElementById('merchantName').value = merchant.merchant_name;
  document.getElementById('merchantType').value = merchant.merchant_type;

  // Store merchant ID in the form
  document.getElementById('addMerchantForm').dataset.merchantId = merchantId;

  // Show the modal
  const modal = new bootstrap.Modal(document.getElementById('addMerchantModal'));
  modal.show();
}

function showAlert(message, type = 'success') {
  const alertDiv = document.createElement('div');
  alertDiv.className = `alert alert-${type} alert-dismissible fade show`;
  alertDiv.innerHTML = `
    ${message}
    <button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>
  `;
  document.querySelector('.container-fluid').insertBefore(alertDiv, document.querySelector('.row'));
  
  // Auto dismiss after 5 seconds
  setTimeout(() => {
    alertDiv.remove();
  }, 5000);
}

async function deleteMerchant(merchantId) {
  if (!confirm('Are you sure you want to deactivate this merchant?')) {
    return;
  }

  try {
    const response = await fetch(`http://localhost:8080/api/merchants/${merchantId}`, {
      method: 'DELETE',
      headers: {
        'accept': 'application/json',
        'Authorization': `Bearer ${getCookie('auth_token')}`,
        'X-User-Id': getCookie('user_id')
      }
    });

    if (response.ok) {
      showAlert('Merchant successfully deactivated!');
      // Reload the merchants table
      const merchantsData = await fetchMerchants();
      if (merchantsData) {
        updateTable(merchantsData);
      }
    } else {
      throw new Error(`Failed to delete merchant: ${response.statusText}`);
    }
  } catch (error) {
    console.error('Error deleting merchant:', error);
    showAlert(error.message, 'danger');
  }
}

// Add pagination control functions
async function nextPage() {
  currentPage++;
  await loadMerchants();
}

async function prevPage() {
  if (currentPage > 1) {
    currentPage--;
    await loadMerchants();
  }
}

// Load merchants when the page loads
// Add Merchant Form Submission
// Update the form submit event listener
document.addEventListener('DOMContentLoaded', function() {
  // Initialize pagination controls
  const pageSizeSelect = document.getElementById('pageSize');
  const prevPageBtn = document.getElementById('prevPage');
  const nextPageBtn = document.getElementById('nextPage');

  // Add event listeners for pagination controls
  if (pageSizeSelect) {
    pageSizeSelect.addEventListener('change', async function() {
      limit = parseInt(this.value);
      currentPage = 1; // Reset to first page when changing page size
      await loadMerchants();
    });
  }

  if (prevPageBtn) {
    prevPageBtn.addEventListener('click', prevPage);
  }

  if (nextPageBtn) {
    nextPageBtn.addEventListener('click', nextPage);
  }

  const addMerchantForm = document.getElementById('addMerchantForm');
  if (addMerchantForm) {
    addMerchantForm.addEventListener('submit', async function(e) {
      e.preventDefault();

      const userId = getCookie('user_id');
      if (!userId) {
        console.error('User ID not found');
        return;
      }

      const formData = {
        merchant_name: document.getElementById('merchantName').value,
        merchant_type: document.getElementById('merchantType').value,
        user_id: userId
      };

      const merchantId = this.dataset.merchantId;
      const isEdit = !!merchantId;

      try {
        const url = isEdit 
          ? `http://localhost:8080/api/merchants/${merchantId}`
          : 'http://localhost:8080/api/merchants';

        const response = await fetch(url, {
          method: isEdit ? 'PUT' : 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${getCookie('auth_token')}`,
            'X-User-Id': userId
          },
          body: JSON.stringify(formData)
        });

        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }

        // Close the modal
        const modal = bootstrap.Modal.getInstance(document.getElementById('addMerchantModal'));
        modal.hide();

        // Reset the form and its state
        addMerchantForm.reset();
        delete addMerchantForm.dataset.merchantId;
        document.getElementById('addMerchantModalLabel').textContent = 'Add New Merchant';
        const submitButton = document.querySelector('#addMerchantForm button[type="submit"]');
        submitButton.textContent = 'Add Merchant';

        // Refresh the merchant table
        await loadMerchants();

        // Show success message
        alert(isEdit ? 'Merchant updated successfully!' : 'Merchant added successfully!');
      } catch (error) {
        console.error('Error:', error);
        alert(isEdit ? 'Failed to update merchant' : 'Failed to add merchant');
      }
    });
  }

  // Initial load of merchants
  loadMerchants();
});
document.addEventListener('DOMContentLoaded', () => {
  loadMerchants();
});
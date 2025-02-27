// Global variables for pagination
let currentPage = 1;
let pageSize = 10;
let selectedMerchantId = null;

// Function to get cookie value
function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
    return null;
  }

// Function to populate merchant dropdown
async function populateMerchantDropdown(merchants) {
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
                fetchProgramRules(merchantId);
            });
        });

    } catch (error) {
        console.error('Error fetching merchants:', error);
    }
}

// Function to fetch program rules for a merchant
async function fetchProgramRules(merchantId = 'all', page = currentPage, limit = pageSize) {
    try {
        const url = merchantId === 'all' 
            ? `/api/program-rules?page=${page}&limit=${limit}`
            : `/api/program-rules/by-merchant/${merchantId}?page=${page}&limit=${limit}`;

        const response = await fetch(url, {
            headers: {
                'Accept': 'application/json',
                'Authorization': `Bearer ${getCookie('auth_token')}`
            }
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        if (!data) {
            throw new Error('No data received from server');
        }

        console.log("Data ", data)

        if (!data.data || data.total_items === 0) {
            showNotification('No program rules found', 'info');
        }

        displayProgramRules(data.data);
        updatePaginationControls(data.total_items, data.current_page, data.per_page, data.total_pages);
    } catch (error) {
        console.error('Error fetching program rules:', error);
        showNotification(`Error fetching program rules: ${error.message}`, 'danger');
    }
}

// Function to display program rules in the table
function displayProgramRules(data) {
    const tableBody = document.getElementById('programsTableBody');
    tableBody.innerHTML = ''; // Clear existing rows

    if (!data || data.length === 0) {
        const emptyRow = document.createElement('tr');
        emptyRow.innerHTML = `
            <td colspan="6" class="text-center">
                <p class="text-sm mb-0">No program rules found</p>
            </td>
        `;
        tableBody.appendChild(emptyRow);
        return;
    }

    data.forEach(program => {
        const row = document.createElement('tr');
        row.innerHTML = `
            <td>
                <div class="d-flex px-3 py-1">
                    <div class="d-flex flex-column justify-content-center">
                        <h6 class="mb-0 text-sm">${program.program_name}</h6>
                        <p class="text-xs text-secondary mb-0">${program.rule_name}</p>
                    </div>
                </div>
            </td>
            <td>
                <p class="text-sm font-weight-bold mb-0">${program.condition_type}</p>
                <p class="text-xs text-secondary mb-0">${program.condition_value}</p>
            </td>
            <td class="align-middle text-center text-sm">
                <span class="badge badge-sm bg-gradient-success">${program.multiplier}x</span>
            </td>
            <td class="align-middle text-center">
                <span class="text-secondary text-sm font-weight-bold">${program.points_awarded}</span>
            </td>
            <td class="align-middle text-center">
                <span class="text-secondary text-xs font-weight-bold">
                    ${formatDate(program.effective_from)} - ${formatDate(program.effective_to)}
                </span>
            </td>
            <td class="align-middle text-center">
                <span class="badge badge-sm bg-gradient-success">Active</span>
            </td>
        `;
        tableBody.appendChild(row);
    });
}

// Function to update pagination controls
function updatePaginationControls(total_items, current_page, per_page, total_pages) {
    // Update page info display
    document.getElementById('currentPage').textContent = current_page || 1;
    document.getElementById('totalPages').textContent = total_pages || 1;

    // Update page size dropdown
    const pageSizeSelect = document.getElementById('pageSize');
    pageSizeSelect.value = per_page;

    // Update button states
    const prevButton = document.getElementById('prevPage');
    const nextButton = document.getElementById('nextPage');

    prevButton.disabled = current_page <= 1;
    nextButton.disabled = current_page >= total_pages;

    // Add event listeners for pagination controls
    prevButton.onclick = () => {
        if (current_page > 1) {
            currentPage = current_page - 1;
            fetchProgramRules(selectedMerchantId, currentPage, per_page);
        }
    };

    nextButton.onclick = () => {
        if (current_page < total_pages) {
            currentPage = current_page + 1;
            fetchProgramRules(selectedMerchantId, currentPage, per_page);
        }
    };

    // Add event listener for page size changes
    pageSizeSelect.onchange = (e) => {
        pageSize = parseInt(e.target.value);
        currentPage = 1; // Reset to first page when changing page size
        fetchProgramRules(selectedMerchantId, currentPage, pageSize);
    };
}

// Function to format date
function formatDate(dateString) {
    if (!dateString) return 'N/A';
    const date = new Date(dateString);
    if (isNaN(date.getTime())) return 'N/A';
    return date.toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric'
    });
}

// Function to show notifications
function showNotification(message, type) {
    // You can implement your notification system here
    console.log(`${type}: ${message}`);
}

// Initialize the page
document.addEventListener('DOMContentLoaded', function() {
        // Initialize merchant dropdown
        populateMerchantDropdown();
});
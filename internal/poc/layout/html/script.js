// --- Bulk Selection Logic ---
function toggleSelectAll() {
    const master = document.getElementById('select-all');
    const checkboxes = document.querySelectorAll('.table-checkbox[name="item"]');
    checkboxes.forEach(cb => cb.checked = master.checked);
}

function selectAll() {
    const master = document.getElementById('select-all');
    const checkboxes = document.querySelectorAll('.table-checkbox[name="item"]');
    master.checked = true;
    checkboxes.forEach(cb => cb.checked = true);
}

function deselectAll() {
    const master = document.getElementById('select-all');
    const checkboxes = document.querySelectorAll('.table-checkbox[name="item"]');
    master.checked = false;
    checkboxes.forEach(cb => cb.checked = false);
}

function invertSelection() {
    const master = document.getElementById('select-all');
    const checkboxes = document.querySelectorAll('.table-checkbox[name="item"]');
    let allChecked = true;
    checkboxes.forEach(cb => {
        cb.checked = !cb.checked;
        if(!cb.checked) allChecked = false;
    });
    master.checked = allChecked;
}

function sortTable(column) {
    const headers = document.querySelectorAll('th.sortable');
    headers.forEach(h => {
        h.classList.remove('asc', 'desc');
        if(h.innerText.toLowerCase().includes(column)) {
            h.classList.add('desc');
        }
    });
}

// --- Filtering Logic ---

function applyFilters() {
    const table = document.querySelector('.data-table');
    const rows = table.querySelectorAll('tbody tr');

    // Collect active filters
    const filters = [];
    const operators = document.querySelectorAll('.filter-operator');
    const inputs = document.querySelectorAll('.filter-input');

    operators.forEach((op, index) => {
        const colName = op.getAttribute('data-col');
        const operator = op.value;
        const value = inputs[index].value.trim();

        if (value !== "") {
            filters.push({
                col: colName,
                op: operator,
                val: value
            });
        }
    });

    // Apply filters to rows
    rows.forEach(row => {
        let isVisible = true;
        const cells = row.querySelectorAll('td');

        // Map column index to data attributes (skipping checkbox col at index 0)
        // Index mapping: 0=checkbox, 1=id, 2=title, 3=email, 4=category, 5=status, 6=created, 7=actions
        const cellMap = {
            'id': 1,
            'title': 2,
            'email': 3,
            'category': 4,
            'status': 5,
            'created': 6
        };

        for (let filter of filters) {
            const cellIndex = cellMap[filter.col];
            if (cellIndex === undefined) continue;

            let cellText = cells[cellIndex].innerText.trim();

            // Special handling for badges (extract text)
            if (filter.col === 'status') {
                const badge = cells[cellIndex].querySelector('.badge');
                if (badge) cellText = badge.innerText.trim();
            }

            // Clean ID (remove #)
            if (filter.col === 'id') {
                cellText = cellText.replace('#', '');
            }

            if (!evaluateCondition(cellText, filter.op, filter.val)) {
                isVisible = false;
                break;
            }
        }

        row.style.display = isVisible ? '' : 'none';
    });
}

function evaluateCondition(cellValue, operator, filterValue) {
    // Convert to lowercase for string comparisons
    const cVal = String(cellValue).toLowerCase();
    const fVal = String(filterValue).toLowerCase();

    switch (operator) {
        case 'eq': return cVal === fVal;
        case 'neq': return cVal !== fVal;
        case 'gt': return parseFloat(cVal) > parseFloat(fVal);
        case 'lt': return parseFloat(cVal) < parseFloat(fVal);
        case 'gte': return parseFloat(cVal) >= parseFloat(fVal);
        case 'lte': return parseFloat(cVal) <= parseFloat(fVal);
        case 'like':
            // SQL LIKE simulation: % becomes .*
            const regexLike = new RegExp('^' + fVal.replace(/%/g, '.*') + '$', 'i');
            return regexLike.test(cVal);
        case 'match':
            // Regex match
            try {
                const regexMatch = new RegExp(fVal, 'i');
                return regexMatch.test(cVal);
            } catch (e) {
                console.error("Invalid Regex:", e);
                return false;
            }
        default: return true;
    }
}

function resetFilters() {
    // Clear inputs
    document.querySelectorAll('.filter-input').forEach(input => input.value = '');
    // Reset operators to default (first option)
    document.querySelectorAll('.filter-operator').forEach(select => select.selectedIndex = 0);
    // Show all rows
    document.querySelectorAll('.data-table tbody tr').forEach(row => row.style.display = '');
}
// Global state
let massageTypes = [];
let selectedService = null;
let selectedDate = null;
let selectedTime = null;
let currentMonth = new Date();
let timeSlots = [];

// DOM elements
const servicesContainer = document.getElementById('services-container');
const calendarSection = document.getElementById('calendar-section');
const selectedServiceInfo = document.getElementById('selected-service-info');
const calendarGrid = document.getElementById('calendar-grid');
const currentMonthElement = document.getElementById('current-month');
const timeSlotsSection = document.getElementById('time-slots-section');
const timeSlotsGrid = document.getElementById('time-slots-grid');
const selectedDateElement = document.getElementById('selected-date');
const selectionSummary = document.getElementById('selection-summary');
const loading = document.getElementById('loading');
const errorMessage = document.getElementById('error-message');
const errorText = document.getElementById('error-text');

// Initialize the application
document.addEventListener('DOMContentLoaded', function() {
    loadMassageTypes();
    setupEventListeners();
});

// Setup event listeners
function setupEventListeners() {
    document.getElementById('prev-month').addEventListener('click', () => {
        currentMonth.setMonth(currentMonth.getMonth() - 1);
        renderCalendar();
    });

    document.getElementById('next-month').addEventListener('click', () => {
        currentMonth.setMonth(currentMonth.getMonth() + 1);
        renderCalendar();
    });
}

// API functions
async function loadMassageTypes() {
    try {
        showLoading(true);
        const response = await fetch('/api/massage-types');
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        massageTypes = await response.json();
        renderServices();
    } catch (error) {
        console.error('Error loading massage types:', error);
        showError('Failed to load massage services. Please try again.');
    } finally {
        showLoading(false);
    }
}

async function loadTimeSlots(date, serviceId) {
    try {
        showLoading(true);
        const response = await fetch(`/api/slots?date=${date}&service_id=${serviceId}`);
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        timeSlots = await response.json();
        renderTimeSlots();
    } catch (error) {
        console.error('Error loading time slots:', error);
        showError('Failed to load available times. Please try again.');
    } finally {
        showLoading(false);
    }
}

// Render functions
function renderServices() {
    servicesContainer.innerHTML = '';
    
    massageTypes.forEach(service => {
        const serviceCard = document.createElement('div');
        serviceCard.className = 'service-card';
        serviceCard.onclick = () => selectService(service);
        
        serviceCard.innerHTML = `
            <h3>${service.name}</h3>
            <div class="service-details">
                <span class="duration">${service.duration} minutes</span>
                <span class="price">€${service.price}</span>
            </div>
        `;
        
        servicesContainer.appendChild(serviceCard);
    });
}

function selectService(service) {
    selectedService = service;
    selectedDate = null;
    selectedTime = null;
    
    // Update UI
    document.querySelectorAll('.service-card').forEach(card => {
        card.classList.remove('selected');
    });
    
    event.currentTarget.classList.add('selected');
    
    // Show calendar section
    calendarSection.style.display = 'block';
    
    // Update selected service info
    selectedServiceInfo.innerHTML = `
        <strong>Selected Service:</strong> ${service.name} 
        (${service.duration} minutes, €${service.price})
    `;
    
    // Hide time slots and summary
    timeSlotsSection.style.display = 'none';
    selectionSummary.style.display = 'none';
    
    // Render calendar
    renderCalendar();
    
    // Scroll to calendar
    calendarSection.scrollIntoView({ behavior: 'smooth' });
}

function renderCalendar() {
    if (!selectedService) return;

    const year = currentMonth.getFullYear();
    const month = currentMonth.getMonth();

    // Update month display
    currentMonthElement.textContent = currentMonth.toLocaleDateString('en-US', {
        month: 'long',
        year: 'numeric'
    });

    // Clear calendar grid
    calendarGrid.innerHTML = '';

    // Add day headers
    const dayHeaders = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];
    dayHeaders.forEach(day => {
        const dayHeader = document.createElement('div');
        dayHeader.textContent = day;
        dayHeader.style.fontWeight = 'bold';
        dayHeader.style.textAlign = 'center';
        dayHeader.style.padding = '10px';
        dayHeader.style.background = '#f5f5f5';
        calendarGrid.appendChild(dayHeader);
    });

    // Get first day of month and number of days
    const firstDay = new Date(year, month, 1);
    const lastDay = new Date(year, month + 1, 0);
    const daysInMonth = lastDay.getDate();
    const startingDayOfWeek = firstDay.getDay();

    // Add empty cells for days before the first day of the month
    for (let i = 0; i < startingDayOfWeek; i++) {
        const emptyDay = document.createElement('div');
        emptyDay.className = 'calendar-day disabled';
        calendarGrid.appendChild(emptyDay);
    }

    // Add days of the month
    const today = new Date();
    for (let day = 1; day <= daysInMonth; day++) {
        const dayElement = document.createElement('div');
        dayElement.className = 'calendar-day';
        dayElement.textContent = day;

        const currentDate = new Date(year, month, day);
        const dateString = currentDate.toISOString().split('T')[0];

        // Disable past dates
        if (currentDate < today.setHours(0, 0, 0, 0)) {
            dayElement.classList.add('disabled');
        } else {
            dayElement.onclick = () => selectDate(dateString);
        }

        // Highlight selected date
        if (selectedDate === dateString) {
            dayElement.classList.add('selected');
        }

        calendarGrid.appendChild(dayElement);
    }
}

function selectDate(dateString) {
    selectedDate = dateString;
    selectedTime = null;

    // Update calendar display
    document.querySelectorAll('.calendar-day').forEach(day => {
        day.classList.remove('selected');
    });

    event.currentTarget.classList.add('selected');

    // Update selected date display
    const date = new Date(dateString);
    selectedDateElement.textContent = date.toLocaleDateString('en-US', {
        weekday: 'long',
        year: 'numeric',
        month: 'long',
        day: 'numeric'
    });

    // Show time slots section
    timeSlotsSection.style.display = 'block';

    // Hide selection summary
    selectionSummary.style.display = 'none';

    // Load time slots for this date
    loadTimeSlots(dateString, selectedService.id);

    // Scroll to time slots
    timeSlotsSection.scrollIntoView({ behavior: 'smooth' });
}

function renderTimeSlots() {
    timeSlotsGrid.innerHTML = '';

    timeSlots.forEach(slot => {
        const slotElement = document.createElement('div');
        slotElement.className = 'time-slot';
        slotElement.textContent = slot.time;

        if (slot.available) {
            slotElement.classList.add('available');
            slotElement.onclick = () => selectTimeSlot(slot);
        } else {
            slotElement.classList.add('booked');
        }

        // Highlight selected time
        if (selectedTime && selectedTime.time === slot.time) {
            slotElement.classList.add('selected');
        }

        timeSlotsGrid.appendChild(slotElement);
    });
}

function selectTimeSlot(slot) {
    selectedTime = slot;

    // Update time slots display
    document.querySelectorAll('.time-slot').forEach(timeSlot => {
        timeSlot.classList.remove('selected');
    });

    event.currentTarget.classList.add('selected');

    // Show selection summary
    updateSelectionSummary();
    selectionSummary.style.display = 'block';

    // Scroll to summary
    selectionSummary.scrollIntoView({ behavior: 'smooth' });
}

function updateSelectionSummary() {
    if (!selectedService || !selectedDate || !selectedTime) return;

    const date = new Date(selectedDate);
    const formattedDate = date.toLocaleDateString('en-US', {
        weekday: 'long',
        year: 'numeric',
        month: 'long',
        day: 'numeric'
    });

    document.getElementById('summary-service').textContent = selectedService.name;
    document.getElementById('summary-date').textContent = formattedDate;
    document.getElementById('summary-time').textContent = selectedTime.time;
    document.getElementById('summary-duration').textContent = selectedService.duration;
    document.getElementById('summary-price').textContent = selectedService.price;
}

// Utility functions
function showLoading(show) {
    loading.style.display = show ? 'flex' : 'none';
}

function showError(message) {
    errorText.textContent = message;
    errorMessage.style.display = 'block';

    // Auto-hide after 5 seconds
    setTimeout(() => {
        hideError();
    }, 5000);
}

function hideError() {
    errorMessage.style.display = 'none';
}

// Booking function (placeholder for now)
function bookAppointment() {
    if (!selectedService || !selectedDate || !selectedTime) {
        showError('Please select a service, date, and time before booking.');
        return;
    }

    // For now, just show a success message
    alert(`Appointment booked successfully!\n\nService: ${selectedService.name}\nDate: ${selectedDate}\nTime: ${selectedTime.time}\nPrice: €${selectedService.price}`);

    // Reset the form
    resetBooking();
}

function resetBooking() {
    selectedService = null;
    selectedDate = null;
    selectedTime = null;

    // Hide sections
    calendarSection.style.display = 'none';
    timeSlotsSection.style.display = 'none';
    selectionSummary.style.display = 'none';

    // Clear selections
    document.querySelectorAll('.service-card').forEach(card => {
        card.classList.remove('selected');
    });

    // Scroll to top
    window.scrollTo({ top: 0, behavior: 'smooth' });
}

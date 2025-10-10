// Global state
let massageTypes = [];
let selectedService = null;
let selectedDate = null;
let selectedTime = null;
let selectedSlotId = null;
let currentMonth = new Date();
let timeSlots = [];

// Story #2 state
let currentReservation = null;
let reservationTimer = null;
let formValidation = {
    name: false,
    email: false,
    phone: false
};

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

// Story #2 DOM elements
const bookingFormSection = document.getElementById('booking-form-section');
const bookingForm = document.getElementById('booking-form');
const reservationTimerElement = document.getElementById('reservation-timer');
const timerDisplay = document.getElementById('timer-display');
const confirmBookingBtn = document.getElementById('confirm-booking-btn');
const clientNameInput = document.getElementById('client-name');
const clientEmailInput = document.getElementById('client-email');
const clientPhoneInput = document.getElementById('client-phone');
const nameError = document.getElementById('name-error');
const emailError = document.getElementById('email-error');
const phoneError = document.getElementById('phone-error');

// Initialize the application
document.addEventListener('DOMContentLoaded', function() {
    loadMassageTypes();
    setupEventListeners();
    setupFormValidation();
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
    selectedSlotId = slot.id; // Store slot ID for reservation

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
    selectedSlotId = null;

    // Hide sections
    calendarSection.style.display = 'none';
    timeSlotsSection.style.display = 'none';
    selectionSummary.style.display = 'none';
    bookingFormSection.style.display = 'none';

    // Clear selections
    document.querySelectorAll('.service-card').forEach(card => {
        card.classList.remove('selected');
    });

    // Scroll to top
    window.scrollTo({ top: 0, behavior: 'smooth' });
}

// Story #2 Functions

// Show booking form and create reservation
async function showBookingForm() {
    if (!selectedSlotId) {
        showError('Please select a time slot first');
        return;
    }

    try {
        showLoading(true);

        // Create temporary reservation
        const response = await fetch('/api/reservations', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                slot_id: selectedSlotId
            })
        });

        if (!response.ok) {
            const errorData = await response.text();
            throw new Error(errorData || 'Failed to create reservation');
        }

        const reservation = await response.json();
        currentReservation = reservation;

        // Update booking summary
        updateBookingSummary();

        // Show booking form
        calendarSection.style.display = 'none';
        bookingFormSection.style.display = 'block';

        // Start countdown timer
        startReservationTimer(reservation.expires_in_seconds);

        // Clear form
        resetForm();

        showLoading(false);

        // Scroll to form
        bookingFormSection.scrollIntoView({ behavior: 'smooth' });

    } catch (error) {
        showLoading(false);
        showError('Failed to reserve time slot: ' + error.message);
    }
}

// Update booking summary in the form
function updateBookingSummary() {
    if (!selectedService || !selectedDate || !selectedTime) return;

    const date = new Date(selectedDate);
    const formattedDate = date.toLocaleDateString('en-US', {
        weekday: 'long',
        year: 'numeric',
        month: 'long',
        day: 'numeric'
    });

    document.getElementById('booking-summary-service').textContent = selectedService.name;
    document.getElementById('booking-summary-date').textContent = formattedDate;
    document.getElementById('booking-summary-time').textContent = selectedTime.time;
    document.getElementById('booking-summary-price').textContent = selectedService.price;
}

// Start reservation countdown timer
function startReservationTimer(seconds) {
    let timeLeft = seconds;

    // Clear any existing timer
    if (reservationTimer) {
        clearInterval(reservationTimer);
    }

    // Update timer display immediately
    updateTimerDisplay(timeLeft);

    reservationTimer = setInterval(() => {
        timeLeft--;
        updateTimerDisplay(timeLeft);

        // Change color when less than 2 minutes
        if (timeLeft <= 120) {
            reservationTimerElement.classList.add('warning');
        }

        // Timer expired
        if (timeLeft <= 0) {
            clearInterval(reservationTimer);
            handleReservationExpired();
        }
    }, 1000);
}

// Update timer display
function updateTimerDisplay(seconds) {
    const minutes = Math.floor(seconds / 60);
    const remainingSeconds = seconds % 60;
    timerDisplay.textContent = `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`;
}

// Handle reservation expiration
function handleReservationExpired() {
    showError('Your reservation has expired. Please select a new time slot.');
    cancelBooking();
}

// Setup form validation
function setupFormValidation() {
    // Name validation
    clientNameInput.addEventListener('input', () => validateName());
    clientNameInput.addEventListener('blur', () => validateName());

    // Email validation
    clientEmailInput.addEventListener('input', () => validateEmail());
    clientEmailInput.addEventListener('blur', () => validateEmail());

    // Phone validation
    clientPhoneInput.addEventListener('input', () => validatePhone());
    clientPhoneInput.addEventListener('blur', () => validatePhone());

    // Form submission
    bookingForm.addEventListener('submit', handleFormSubmit);
}

// Validation functions
function validateName() {
    const name = clientNameInput.value.trim();
    let isValid = true;
    let errorMessage = '';

    if (!name) {
        isValid = false;
        errorMessage = 'Name is required';
    } else if (name.length < 2) {
        isValid = false;
        errorMessage = 'Name must be at least 2 characters';
    } else if (!/^[a-zA-Z\s]+$/.test(name)) {
        isValid = false;
        errorMessage = 'Name should contain only letters and spaces';
    }

    formValidation.name = isValid;
    nameError.textContent = errorMessage;
    clientNameInput.classList.toggle('error', !isValid);
    clientNameInput.classList.toggle('valid', isValid && name.length > 0);

    updateSubmitButton();
    return isValid;
}

function validateEmail() {
    const email = clientEmailInput.value.trim();
    let isValid = true;
    let errorMessage = '';

    if (!email) {
        isValid = false;
        errorMessage = 'Email is required';
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
        isValid = false;
        errorMessage = 'Please enter a valid email';
    }

    formValidation.email = isValid;
    emailError.textContent = errorMessage;
    clientEmailInput.classList.toggle('error', !isValid);
    clientEmailInput.classList.toggle('valid', isValid && email.length > 0);

    updateSubmitButton();
    return isValid;
}

function validatePhone() {
    const phone = clientPhoneInput.value.trim();
    let isValid = true;
    let errorMessage = '';

    if (!phone) {
        isValid = false;
        errorMessage = 'Phone is required';
    } else if (!/^\+?[\d\s\-()]{8,}$/.test(phone)) {
        isValid = false;
        errorMessage = 'Please enter a valid phone number';
    }

    formValidation.phone = isValid;
    phoneError.textContent = errorMessage;
    clientPhoneInput.classList.toggle('error', !isValid);
    clientPhoneInput.classList.toggle('valid', isValid && phone.length > 0);

    updateSubmitButton();
    return isValid;
}

// Update submit button state
function updateSubmitButton() {
    const allValid = formValidation.name && formValidation.email && formValidation.phone;
    confirmBookingBtn.disabled = !allValid;
}

// Handle form submission
async function handleFormSubmit(event) {
    event.preventDefault();

    // Validate all fields
    const nameValid = validateName();
    const emailValid = validateEmail();
    const phoneValid = validatePhone();

    if (!nameValid || !emailValid || !phoneValid) {
        return;
    }

    if (!currentReservation) {
        showError('No active reservation found');
        return;
    }

    try {
        confirmBookingBtn.disabled = true;
        confirmBookingBtn.textContent = 'Creating Booking...';

        const bookingData = {
            reservation_id: currentReservation.reservation_id,
            client_name: clientNameInput.value.trim(),
            email: clientEmailInput.value.trim(),
            phone: clientPhoneInput.value.trim(),
            service_id: selectedService.id,
            date: selectedDate,
            time_slot: selectedTime.time
        };

        const response = await fetch('/api/bookings', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(bookingData)
        });

        if (!response.ok) {
            const errorData = await response.text();
            throw new Error(errorData || 'Failed to create booking');
        }

        const booking = await response.json();

        // Clear timer
        if (reservationTimer) {
            clearInterval(reservationTimer);
        }

        // Redirect to confirmation page with booking ID
        window.location.href = `/confirmation.html?id=${booking.id}`;

        // Note: No need to reset state as we're navigating away

    } catch (error) {
        showError('Failed to create booking: ' + error.message);
        confirmBookingBtn.disabled = false;
        confirmBookingBtn.textContent = 'Confirm Booking';
    }
}

// Cancel booking and return to calendar
function cancelBooking() {
    // Cancel reservation if exists
    if (currentReservation) {
        cancelReservation();
    }

    // Clear timer
    if (reservationTimer) {
        clearInterval(reservationTimer);
    }

    // Hide form and show calendar
    bookingFormSection.style.display = 'none';
    calendarSection.style.display = 'block';

    // Reset form
    resetForm();
}

// Cancel reservation on server
async function cancelReservation() {
    if (!currentReservation) return;

    try {
        await fetch(`/api/reservations/${currentReservation.reservation_id}`, {
            method: 'DELETE'
        });
    } catch (error) {
        console.error('Failed to cancel reservation:', error);
    }

    currentReservation = null;
}

// Reset form
function resetForm() {
    bookingForm.reset();
    formValidation = { name: false, email: false, phone: false };

    // Clear error messages
    nameError.textContent = '';
    emailError.textContent = '';
    phoneError.textContent = '';

    // Clear validation classes
    clientNameInput.classList.remove('error', 'valid');
    clientEmailInput.classList.remove('error', 'valid');
    clientPhoneInput.classList.remove('error', 'valid');

    // Reset button
    confirmBookingBtn.disabled = true;
    confirmBookingBtn.textContent = 'Confirm Booking';

    // Reset timer display
    reservationTimerElement.classList.remove('warning');
}

// Reset application state
function resetApplicationState() {
    selectedService = null;
    selectedDate = null;
    selectedTime = null;
    selectedSlotId = null;
    currentReservation = null;

    // Hide all sections
    calendarSection.style.display = 'none';
    bookingFormSection.style.display = 'none';
    selectionSummary.style.display = 'none';
    timeSlotsSection.style.display = 'none';

    // Clear selections
    document.querySelectorAll('.service-card.selected').forEach(card => {
        card.classList.remove('selected');
    });

    // Scroll to top
    window.scrollTo({ top: 0, behavior: 'smooth' });
}

// Show success message
function showSuccess(message) {
    // Create success message element (reuse error message styling)
    const successDiv = document.createElement('div');
    successDiv.className = 'success-message';
    successDiv.innerHTML = `
        <p>${message}</p>
        <button onclick="this.parentElement.remove()">Close</button>
    `;
    successDiv.style.cssText = `
        position: fixed;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        background: #4CAF50;
        color: white;
        padding: 20px;
        border-radius: 8px;
        box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
        z-index: 1000;
        text-align: center;
    `;

    document.body.appendChild(successDiv);

    // Auto-remove after 5 seconds
    setTimeout(() => {
        if (successDiv.parentElement) {
            successDiv.remove();
        }
    }, 5000);
}

// Placeholder function for backward compatibility
function bookAppointment() {
    showBookingForm();
}

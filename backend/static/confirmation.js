// Confirmation page JavaScript

// DOM elements
const loadingSection = document.getElementById('loading-section');
const errorSection = document.getElementById('error-section');
const successSection = document.getElementById('success-section');

// Booking detail elements
const bookingReference = document.getElementById('booking-reference');
const serviceName = document.getElementById('service-name');
const serviceDuration = document.getElementById('service-duration');
const servicePrice = document.getElementById('service-price');
const bookingDate = document.getElementById('booking-date');
const bookingTime = document.getElementById('booking-time');
const clientName = document.getElementById('client-name');
const clientEmail = document.getElementById('client-email');
const clientPhone = document.getElementById('client-phone');

// Initialize page when DOM is loaded
document.addEventListener('DOMContentLoaded', function() {
    loadBookingDetails();
});

// Load booking details from API
async function loadBookingDetails() {
    try {
        // Get booking ID from URL parameters
        const urlParams = new URLSearchParams(window.location.search);
        const bookingId = urlParams.get('id');
        
        if (!bookingId) {
            showError();
            return;
        }

        // Fetch booking details from API
        const response = await fetch(`/api/bookings/${bookingId}`);
        
        if (!response.ok) {
            if (response.status === 404) {
                showError();
                return;
            }
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const booking = await response.json();
        displayBookingDetails(booking);
        
    } catch (error) {
        console.error('Error loading booking details:', error);
        showError();
    }
}

// Display booking details in the UI
function displayBookingDetails(booking) {
    // Hide loading, show success
    loadingSection.style.display = 'none';
    successSection.style.display = 'block';

    // Populate booking details
    bookingReference.textContent = booking.reference;
    serviceName.textContent = booking.service_name;
    serviceDuration.textContent = `${booking.duration} minutes`;
    servicePrice.textContent = `€${booking.price.toFixed(2)}`;
    
    // Format and display date
    const formattedDate = formatDate(booking.date);
    bookingDate.textContent = formattedDate;
    
    bookingTime.textContent = booking.time_slot;
    clientName.textContent = booking.client_name;
    clientEmail.textContent = booking.email;
    clientPhone.textContent = booking.phone;

    // Update page title with reference
    document.title = `Booking Confirmation ${booking.reference} - Massage Booking`;
}

// Show error state
function showError() {
    loadingSection.style.display = 'none';
    errorSection.style.display = 'block';
}

// Format date for display
function formatDate(dateString) {
    try {
        const date = new Date(dateString);
        const options = { 
            weekday: 'long', 
            year: 'numeric', 
            month: 'long', 
            day: 'numeric' 
        };
        return date.toLocaleDateString('en-US', options);
    } catch (error) {
        console.error('Error formatting date:', error);
        return dateString; // Fallback to original string
    }
}

// Navigate back to home page
function goHome() {
    window.location.href = '/';
}

// Handle browser back button
window.addEventListener('popstate', function(event) {
    // If user navigates back, go to home page
    goHome();
});

// Add some visual feedback for the success animation
function animateSuccess() {
    const successIcon = document.querySelector('.success-icon');
    if (successIcon) {
        successIcon.style.animation = 'bounce 0.6s ease-in-out';
    }
}

// Call animation when success section is shown
const observer = new MutationObserver(function(mutations) {
    mutations.forEach(function(mutation) {
        if (mutation.type === 'attributes' && mutation.attributeName === 'style') {
            const target = mutation.target;
            if (target.id === 'success-section' && target.style.display === 'block') {
                setTimeout(animateSuccess, 100);
            }
        }
    });
});

// Start observing
observer.observe(successSection, { attributes: true });

// Add CSS animation for bounce effect
const style = document.createElement('style');
style.textContent = `
    @keyframes bounce {
        0%, 20%, 60%, 100% {
            transform: translateY(0);
        }
        40% {
            transform: translateY(-10px);
        }
        80% {
            transform: translateY(-5px);
        }
    }
`;
document.head.appendChild(style);

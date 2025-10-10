package email

import (
	"fmt"
	"strings"
	"time"

	"massage-booking/backend/models"
)

// RenderEmailTemplate generates HTML email content for booking confirmation
func RenderEmailTemplate(booking *models.BookingDetail) string {
	// Format date for display
	date, err := time.Parse("2006-01-02", booking.Date)
	var formattedDate string
	if err != nil {
		formattedDate = booking.Date
	} else {
		formattedDate = date.Format("Monday, January 2, 2006")
	}

	template := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Booking Confirmation</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f4f4f4;
        }
        .container {
            background: white;
            padding: 30px;
            border-radius: 10px;
            box-shadow: 0 0 10px rgba(0,0,0,0.1);
        }
        .header {
            text-align: center;
            border-bottom: 2px solid #4CAF50;
            padding-bottom: 20px;
            margin-bottom: 30px;
        }
        .header h1 {
            color: #4CAF50;
            margin: 0;
        }
        .booking-details {
            background: #f8f9fa;
            padding: 20px;
            border-radius: 8px;
            margin: 20px 0;
        }
        .detail-row {
            display: flex;
            justify-content: space-between;
            margin-bottom: 10px;
            padding: 5px 0;
            border-bottom: 1px solid #eee;
        }
        .detail-label {
            font-weight: bold;
            color: #555;
        }
        .detail-value {
            color: #333;
        }
        .reference {
            background: #e3f2fd;
            padding: 15px;
            border-radius: 8px;
            text-align: center;
            margin: 20px 0;
            border-left: 4px solid #2196f3;
        }
        .reference-number {
            font-size: 18px;
            font-weight: bold;
            color: #1976d2;
        }
        .footer {
            text-align: center;
            margin-top: 30px;
            padding-top: 20px;
            border-top: 1px solid #eee;
            color: #666;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Booking Confirmation</h1>
            <p>Your massage appointment has been confirmed!</p>
        </div>

        <div class="reference">
            <p>Booking Reference Number:</p>
            <div class="reference-number">{{.Reference}}</div>
        </div>

        <div class="booking-details">
            <h3>Appointment Details</h3>
            <div class="detail-row">
                <span class="detail-label">Service:</span>
                <span class="detail-value">{{.ServiceName}}</span>
            </div>
            <div class="detail-row">
                <span class="detail-label">Duration:</span>
                <span class="detail-value">{{.Duration}} minutes</span>
            </div>
            <div class="detail-row">
                <span class="detail-label">Price:</span>
                <span class="detail-value">€{{.Price}}</span>
            </div>
            <div class="detail-row">
                <span class="detail-label">Date:</span>
                <span class="detail-value">{{.FormattedDate}}</span>
            </div>
            <div class="detail-row">
                <span class="detail-label">Time:</span>
                <span class="detail-value">{{.TimeSlot}}</span>
            </div>
        </div>

        <div class="booking-details">
            <h3>Customer Information</h3>
            <div class="detail-row">
                <span class="detail-label">Name:</span>
                <span class="detail-value">{{.ClientName}}</span>
            </div>
            <div class="detail-row">
                <span class="detail-label">Email:</span>
                <span class="detail-value">{{.Email}}</span>
            </div>
            <div class="detail-row">
                <span class="detail-label">Phone:</span>
                <span class="detail-value">{{.Phone}}</span>
            </div>
        </div>

        <div class="footer">
            <p>We look forward to seeing you!</p>
            <p><strong>Massage Booking Team</strong></p>
            <p style="font-size: 12px; color: #999;">
                Please save this email for your records. If you need to make any changes, 
                please contact us with your booking reference number.
            </p>
        </div>
    </div>
</body>
</html>`

	// Replace placeholders
	content := strings.ReplaceAll(template, "{{.Reference}}", booking.Reference)
	content = strings.ReplaceAll(content, "{{.ServiceName}}", booking.ServiceName)
	content = strings.ReplaceAll(content, "{{.Duration}}", fmt.Sprintf("%d", booking.Duration))
	content = strings.ReplaceAll(content, "{{.Price}}", fmt.Sprintf("%.2f", booking.Price))
	content = strings.ReplaceAll(content, "{{.FormattedDate}}", formattedDate)
	content = strings.ReplaceAll(content, "{{.TimeSlot}}", booking.TimeSlot)
	content = strings.ReplaceAll(content, "{{.ClientName}}", booking.ClientName)
	content = strings.ReplaceAll(content, "{{.Email}}", booking.Email)
	content = strings.ReplaceAll(content, "{{.Phone}}", booking.Phone)

	return content
}

// GetEmailSubject generates email subject line
func GetEmailSubject(booking *models.BookingDetail) string {
	return fmt.Sprintf("Booking Confirmation - %s", booking.Reference)
}

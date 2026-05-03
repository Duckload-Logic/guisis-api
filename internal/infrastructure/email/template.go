package email

func OTP_TEMPLATE(otp string) string {
	return `
	<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>One-Time Password</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background-color: #f4f4f9;
            margin: 0;
            padding: 0;
            display: flex;
            justify-content: center;
            align-items: center;
            height: 100vh;
        }
        .container {
            position: relative;
            background-color: #ffffff;
            padding: 40px;
            border-radius: 12px;
            box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
            max-width: 400px;
            width: 100%;
            text-align: center;
            overflow: hidden; /* Ensures the highlight bar follows the container's corners */
        }
        .logo {
			display: flex;
			align-items: center;
			justify-content: center;
			gap: 10px;
            font-size: 28px;
            font-weight: 700;
            color: #2c3e50;
            margin-bottom: 20px;
        }
        .title {
            font-size: 24px;
            font-weight: 600;
            color: #630b0bff;
            margin-bottom: 10px;
        }
        .subtitle {
            color: #6c757d;
            margin-bottom: 30px;
            font-size: 14px;
        }
        .otp-box {
            background-color: #e9ecef;
            padding: 20px;
            border-radius: 8px;
            margin: 20px 0;
            display: inline-block;
        }
        .otp-code {
            font-size: 36px;
            font-weight: 700;
            color: #2c3e50;
            letter-spacing: 4px;
        }
        .info-text {
            color: #6c757d;
            font-size: 13px;
            margin-bottom: 20px;
        }
        .footer {
            margin-top: 30px;
            padding-top: 20px;
            border-top: 1px solid #e9ecef;
            color: #95a5a6;
            font-size: 12px;
        }
        .footer a {
            color: #3498db;
            text-decoration: none;
            font-weight: 500;
        }
    </style>
</head>
<body>
    <div class="container">
		<div style="position: absolute; top: 0; left: 0; width: 100%; height: 8px; background-color: #630b0b;"></div>
		<div class="header">
			<div class="logo"><img src="https://guisis.dllbsit2027.com/logo.svg" width="50" height="50"> GuiSIS</div>
		</div>
        <div class="title">One-Time Password</div>
        <div class="subtitle">Please use the code below to verify your identity</div>

        <div class="otp-box">
            <div class="otp-code">` + otp + `</div>
        </div>

        <div class="info-text">
            This code will expire in 5 minutes
        </div>

        <div class="footer">
            If you did not request this code, please ignore this email.
            <br>
            For security reasons, never share your OTP with anyone.
        </div>
    </div>
</body>
</html>
`
}

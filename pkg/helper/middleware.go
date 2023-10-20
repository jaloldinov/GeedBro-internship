package helper

/*
func StartMiddleware(c *gin.Context) {
	// before request
	fmt.Printf("%s request start path: %s time %v\n", c.Request.Method, c.Request.URL.Path, time.Now())
	c.Next()
}
func EndMiddleware(c *gin.Context) {
	c.Next()
	// after request
	fmt.Printf("%s request end path: %s time %v\n", c.Request.Method, c.Request.URL.Path, time.Now())
}
func LoggerAllInOne(c *gin.Context) {
	// before request
	t := time.Now()
	c.Next()
	// after request
	latency := time.Since(t)

	// access the status we are sending
	status := c.Writer.Status()

	fmt.Printf("Completed %s %s with status code %d in %v\n", c.Request.Method, c.Request.URL.Path, status, latency)
}
func Logger(c *gin.Context) {
	// before request
	beforeRequest(c)
	c.Next()
	// after request
	afterRequest(c)
}
func beforeRequest(c *gin.Context) {
	// before request
	t := time.Now()
	c.Set("start", t)
	c.Next()
}
func afterRequest(c *gin.Context) {
	// Get the start time from the request context
	startTime, exists := c.Get("start")
	if !exists {
		startTime = time.Now()
	}

	// Calculate the request duration
	duration := time.Since(startTime.(time.Time))

	// Log the request completion time and duration
	fmt.Printf("Completed %s %s in %v\n", c.Request.Method, c.Request.URL.Path, duration)
}

// AuthMiddleWare is a middleware function for authentication
func AuthMiddleWare(c *gin.Context) {
	token := c.GetHeader("Authorization")

	if token == "" {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    "UNAUTHORIZED!",
			"message": "Token not found...",
		})
		c.Abort()
		return
	}

	staffInfo, err := ParseClaims(token, config.JWTSecretKey)
	if err != nil {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    "INVALID TOKEN!",
			"message": "Provided token is not valid...",
		})
		c.Abort()
		return
	}

	c.Set("staff_info", staffInfo)
	c.Next()
}

func ValidatePasswordMiddleware(c *gin.Context) {
	var requestBody struct {
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body (password)"})
		c.Abort()
		return
	}

	password := requestBody.Password

	if !IsValidPassword(password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password. Password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, and one digit."})
		c.Abort()
		return
	}

	// Continue with other middleware or the final endpoint handler
	c.Next()
}

// login(username) middleware
func ValidateLoginMiddleware(c *gin.Context) {
	var login models.Login

	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body (login)"})
		c.Abort()
		return
	}

	if !IsValidLogin(login.Login) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username. Username must start with a letter, contain only alphanumeric characters and underscores, and be between 6 to 30 characters long."})
		c.Abort()
		return
	}

	c.Next()
}

// phoneNumber middleware
func ValidatePhoneNumberMiddleware(c *gin.Context) {
	var req models.PhoneNumberRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body (number)"})
		c.Abort()
		return
	}

	if !IsValidPhone(req.PhoneNumber) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid phone number. Phone number must start with '+998'"})
		c.Abort()
		return
	}

	c.Next()
}
*/

package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

var tpl = template.Must(template.New("index").Parse(indexHTML))

type UserInfo struct {
	IP           string   `json:"ip"`
	UserAgent    string   `json:"user_agent"`
	Languages    []string `json:"languages"`
	Platform     string   `json:"platform"`
	ScreenWidth  int      `json:"screen_width"`
	ScreenHeight int      `json:"screen_height"`
	ColorDepth   int      `json:"color_depth"`
	Timezone     string   `json:"timezone"`
	Cookies      string   `json:"cookies_enabled"`
	Online       bool     `json:"online"`
	Referrer     string   `json:"referrer"`
}

var collectedData []UserInfo

func main() {
	http.HandleFunc("/", loadingPage)
	http.HandleFunc("/next", redirectPage)
	http.HandleFunc("/collect", collectDataHandler)
	http.HandleFunc("/data", showCollectedData)
	http.HandleFunc("/clear", clearDataHandler)

	fmt.Println("Server starting on http://localhost:8080")
	fmt.Println("Educational demo - this collects user data for learning purposes")
	http.ListenAndServe(":8080", nil)
}

func loadingPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tpl.ExecuteTemplate(w, "index", nil)
}

func redirectPage(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://www.youtube.com/watch?v=dQw4w9WgXcQ", http.StatusSeeOther)
}

func collectDataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		Languages    []string `json:"languages"`
		Platform     string   `json:"platform"`
		ScreenWidth  int      `json:"screenWidth"`
		ScreenHeight int      `json:"screenHeight"`
		ColorDepth   int      `json:"colorDepth"`
		Timezone     string   `json:"timezone"`
		Cookies      string   `json:"cookies"`
		Online       bool     `json:"online"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get additional server-side info
	ip := strings.Split(r.RemoteAddr, ":")[0]
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		ip = forwarded
	}

	userInfo := UserInfo{
		IP:           ip,
		UserAgent:    r.UserAgent(),
		Languages:    data.Languages,
		Platform:     data.Platform,
		ScreenWidth:  data.ScreenWidth,
		ScreenHeight: data.ScreenHeight,
		ColorDepth:   data.ColorDepth,
		Timezone:     data.Timezone,
		Cookies:      data.Cookies,
		Online:       data.Online,
		Referrer:     r.Referer(),
	}

	collectedData = append(collectedData, userInfo)

	// Log to console (for demonstration)
	fmt.Printf("\n=== COLLECTED USER DATA ===\n")
	fmt.Printf("IP: %s\n", userInfo.IP)
	fmt.Printf("User Agent: %s\n", userInfo.UserAgent)
	fmt.Printf("Languages: %v\n", userInfo.Languages)
	fmt.Printf("Platform: %s\n", userInfo.Platform)
	fmt.Printf("Screen: %dx%d\n", userInfo.ScreenWidth, userInfo.ScreenHeight)
	fmt.Printf("Referrer: %s\n", userInfo.Referrer)
	fmt.Printf("===========================\n")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Data collected (for educational purposes)"))
}

func showCollectedData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(collectedData)
}

func clearDataHandler(w http.ResponseWriter, r *http.Request) {
	collectedData = []UserInfo{}
	w.Write([]byte("Data cleared"))
}

const indexHTML = `
<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>Loading... | Educational Demo</title>
<style>
  body {
    background: white;
    margin: 0;
    height: 100vh;
    display: flex;
    justify-content: center;
    align-items: center;
    font-family: Arial, sans-serif;
  }
  .loading-text {
    color: #333;
    font-size: 18px;
  }
  .warning {
    position: fixed;
    bottom: 20px;
    left: 20px;
    background: #ff6b6b;
    color: white;
    padding: 15px;
    border-radius: 5px;
    font-size: 14px;
    max-width: 300px;
    display: none;
  }
</style>
</head>
<body>
  <div class="loading-text">Loading content...</div>
  <div class="warning" id="warning">
    ⚠️ This educational demo collected information about your device. 
    Check the developer console for details.
  </div>

<script>
  // Collect user data
  const userData = {
    languages: navigator.languages || [navigator.language],
    platform: navigator.platform,
    screenWidth: screen.width,
    screenHeight: screen.height,
    colorDepth: screen.colorDepth,
    timezone: Intl.DateTimeFormat().resolvedOptions().timeZone,
    cookies: navigator.cookieEnabled ? "Enabled" : "Disabled",
    online: navigator.onLine
  };

  // Send collected data to server
  fetch('/collect', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(userData)
  })
  .then(response => {
    // Show educational warning
    document.getElementById('warning').style.display = 'block';
    
    // Display collected data in console for educational purposes
    console.log('🔍 EDUCATIONAL DEMO: Information collected about your device:');
    console.log('🌐 Languages:', userData.languages);
    console.log('💻 Platform:', userData.platform);
    console.log('📱 Screen:', userData.screenWidth + 'x' + userData.screenHeight);
    console.log('🎨 Color Depth:', userData.colorDepth + ' bits');
    console.log('⏰ Timezone:', userData.timezone);
    console.log('🍪 Cookies:', userData.cookies);
    console.log('📶 Online:', userData.online);
    console.log('-----------------------------------');
    console.log('This is how some websites collect data without explicit consent.');
    console.log('Always check what information websites are accessing!');
    
    // Wait 3 seconds then redirect
    setTimeout(function() {
      window.location.href = "/next";
    }, 3000);
  })
  .catch(error => {
    console.error('Error:', error);
    window.location.href = "/next";
  });
</script>
</body>
</html>
`

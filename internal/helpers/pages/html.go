package pages

/*
IndexPage renders the html content for the index page.
*/
const IndexPage = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>OAuth Login</title>
</head>
<body>
    <h1>Welcome</h1>
    <p>Please sign in to continue</p>
    
    <div>
        <a href="/login-gl">Login with Google</a>
    </div>
    <div>
        <a href="/login-gh">Login with GitHub</a>
    </div>
</body>
</html>`

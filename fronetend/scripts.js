function showSignup() {
    document.getElementById('login-section').style.display = 'none';
    document.getElementById('signup-section').style.display = 'block';
}

function showLogin() {
    document.getElementById('signup-section').style.display = 'none';
    document.getElementById('login-section').style.display = 'block';
}

let userId = null; // Variable to store user_id

function login() {
    const email = document.getElementById('login-email').value;
    const password = document.getElementById('login-password').value;

    fetch('http://localhost:8080/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ email, password })
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        return response.json();
    })
    .then(data => {
        if (data.message === "Login successful") {
            userId = data.user_id; // Store user_id for later use
            document.getElementById('login-section').style.display = 'none';
            document.getElementById('home-section').style.display = 'block';
            fetchProposals();
        } else {
            alert('Login failed: ' + data.error || 'Invalid credentials');
        }
    })
    .catch(error => {
        console.error('Error:', error);
        alert('Login failed: ' + error.message); // Display network error
    });
}

function signup() {
    const name = document.getElementById('signup-name').value;
    const designation = document.getElementById('signup-designation').value;
    const experience = parseInt(document.getElementById('signup-experience').value, 10); // Convert to integer
    const skills = document.getElementById('signup-skills').value;
    const email = document.getElementById('signup-email').value;
    const password = document.getElementById('signup-password').value;

    // Ensure experience is a number before sending
    if (isNaN(experience)) {
        alert('Experience must be a valid number.');
        return;
    }

    fetch('http://localhost:8080/signup', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ name, designation, experience, skills, email, password })
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        return response.json();
    })
    .then(data => {
        if (data.message === "Signup successful") {
            alert('Signup successful! You can now log in.');
            showLogin(); // Call your function to show the login form
        } else {
            alert('Signup failed: ' + data.error || 'Unexpected error');
        }
    })
    .catch(error => {
        console.error('Error:', error);
        alert('Signup failed: ' + error.message); // Display network error
    });
}


function fetchProposals() {
    fetch('http://localhost:8080/proposals')
    .then(response => response.json())
    .then(data => {
        const proposalsList = document.getElementById('proposals-list');
        proposalsList.innerHTML = ''; // Clear previous proposals

        data.forEach(proposal => {
            const div = document.createElement('a');
            div.className = 'list-group-item list-group-item-action';
            div.innerHTML = `<strong>${proposal.title}</strong> - ${proposal.description}`;
            div.onclick = () => fetchSubtasks(proposal.id, proposal.title);
            proposalsList.appendChild(div);
        });
    })
    .catch(error => console.error('Error:', error));
}

function fetchSubtasks(proposalID, title) {
    document.getElementById('home-section').style.display = 'none';
    document.getElementById('subtasks-section').style.display = 'block';
    document.getElementById('proposal-title').innerText = title;

    fetch(`http://localhost:8080/proposals/${proposalID}/subtasks`)
    .then(response => response.json())
    .then(data => {
        const subtasksList = document.getElementById('subtasks-list');
        subtasksList.innerHTML = ''; // Clear previous subtaskss

        data.forEach(subtask => {
            const div = document.createElement('div');
            div.className = 'list-group-item';
            div.innerHTML = `<strong>${subtask.topic_name}</strong> - ${subtask.description} (Due: ${subtask.submission_date})`;
            subtasksList.appendChild(div);
        });
    })
    .catch(error => console.error('Error:', error));
}

function goBack() {
    document.getElementById('subtasks-section').style.display = 'none';
    document.getElementById('home-section').style.display = 'block';
}


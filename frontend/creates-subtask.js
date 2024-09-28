function openSubtaskForm() {
    $('#subtaskModal').modal('show');
}

document.getElementById('subtaskForm').onsubmit = function (event) {
    event.preventDefault();

    const proposalID = parseInt(new URLSearchParams(window.location.search).get('proposal_id'), 10);
    const topicName = document.getElementById('topic_name').value;
    const description = document.getElementById('description').value;
    const skillsRequired = document.getElementById('skills_required').value;
    const submissionDate = document.getElementById('submission_date').value;

    const newSubtask = {
        proposal_id: proposalID,
        topic_name: topicName,
        description: description,
        skills_required: skillsRequired,
        submission_date: submissionDate,
    };

    console.log(newSubtask)

    fetch('http://localhost:8080/subtask', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(newSubtask),
    })
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok: ' + response.statusText);
            }
            alert('Subtask created successfully!');
            $('#subtaskModal').modal('hide'); // Hide the modal
            location.reload(); // Refresh the page to show the new subtask
        })
        .catch(error => {
            console.error('Error:', error);
            alert('Failed to create subtask: ' + error.message);
        });
};
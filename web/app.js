const statusFilter = document.getElementById("status-filter");
statusFilter.addEventListener("change", async () => {
    await loadTasks();
});

function createTaskElement(task) {
    const taskElement = document.createElement("div");

    taskElement.textContent =
        `${task.ID}: ${task.Description} [${task.Status}]`;

    const editButton = document.createElement("button");
    editButton.textContent = "Edit";

    editButton.addEventListener("click", async () => {
        const newDescription = prompt(
            "Введите новое описание:",
            task.Description
        );

        if (newDescription === null) {
            return;
        }

        await fetch(`/api/tasks/${task.ID}`, {
            method: "PUT",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                description: newDescription,
            }),
        });

        await loadTasks();
    });

    taskElement.appendChild(editButton);

    addStatusButton("todo", "in-progress", task, "In Progress", taskElement)
    addStatusButton("in-progress", "done", task, "Done", taskElement)

    const deleteButton = document.createElement("button");
    deleteButton.textContent = "Delete";

    deleteButton.addEventListener("click", async () => {
        await fetch(`/api/tasks/${task.ID}`, {
            method: "DELETE",
        });

        await loadTasks();
    });

    taskElement.appendChild(deleteButton);

    return taskElement;
}

async function loadTasks() {

    const status = statusFilter.value;

    const response = await fetch(`/api/tasks?status=${status}`);

    const tasks = await response.json();

    const tasksContainer = document.getElementById("tasks");
    tasksContainer.innerHTML = "";

    for (const task of tasks) {
        const taskElement = createTaskElement(task);
        tasksContainer.appendChild(taskElement);
    }
}

const form = document.getElementById("create-task-form");
const descriptionInput = document.getElementById("description");

form.addEventListener("submit", async (event) => {
    event.preventDefault();

    const response = await fetch("/api/tasks", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            description: descriptionInput.value,
        }),
    });

    if (!response.ok) {
        console.error("Failed to create task");
        return;
    }

    descriptionInput.value = "";

    await loadTasks();
});

loadTasks();


function addStatusButton(from, to, task, ButtonName, taskElement) {
    if (task.Status === from) {
        const button = document.createElement("button");
        button.textContent = ButtonName;

        button.addEventListener("click", async () => {
            await fetch(`/api/tasks/${task.ID}/${to}`, {
                method: "POST",
            });

            await loadTasks();
        });

        taskElement.appendChild(button);
    }
}


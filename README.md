# Task Manager REST API

Simple task manager RESTful API service, implementing CRUD task entity, repetition rules and search.
<br>

### Endpoints📞

| URL                 | HTTP      | Action                        |
| ----------------    | --------- | ----------------------------- |
| '/'                 | GET       | get web-interface             |
| '/api/task?id={id}' | GET       | get task by {id}              |
| '/api/task'         | POST      | create task from request body |
| '/api/task'         | PUT       | update task                   |
| '/api/task?id={id}' | DELETE    | delete task by {id}           |
| '/api/tasks'        | GET       | get last 50 tasks             |
| '/api/nextdate'     | GET       | get next date for the task    |
| '/api/task/done?id={id}' | POST | completе the task by {id}     |

#
### Database table🔖

| Column     | Type         | Description                         |
| ---------- | ------------ | ----------------------------------- |
| id         | INT          | unique primary key                  |
| date       | VARCHAR(8)   | task completion date                |
| title      | TEXT         | task name                           |
| comment    | TEXT         | additional task comment             |
| repeat     | VARCHAR(128) | task repetition rule                |

#
### Run app🚀
<br>

1. **Build docker image:**
```bash
docker build -t todo_app .
```
2. **Run docker container:**
```bash
docker run -d -p 7540:7540 todo_app
```
3. **Use API or Web-interface:**
```bash
http://localhost:7540/
```

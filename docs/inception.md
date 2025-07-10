# Idea and Vision

Pain: I need a buffer to store my time logs before syncing them with ComindWork at the end of the day.
For this, I use Clockify with each time log being of special format:

1. Company: `<my company name>`
1. Task: `<project name>/TASK<task number>`
1. What have you worked on?: `<name of time log>`

And this structure perfectly maps to the ComindWork time log format:

1. Project: `<project name>`
1. Task: `<task number>`
1. What have you worked on?: `<name of time log>`

But doing this on my own every day takes some time, which I would rather automate.

Therefore, the goal of this project is to develop an app to synchronize time logs between Clockify and ComindWork.

## Desired features

1. **Main idea**: synchronize time logs between Clockify and ComindWork.
1. **Web app**: to set up configurations and control the synchronization process.
1. **Visualization**: make visualizations of my time logs (e.g. by day, by project, by task, etc.).
1. _Preferably_:

    1. **Real-time synchronization**: instead of pushing a button, synchronize time logs in real time.
    1. **Webhook**: instead of polling the Clockify API, use a webhook to get notified about new time logs.
    1. **AI analysis and recommendations**: analyze my time logs and provide recommendations for better time management.

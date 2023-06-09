# the organizations table
# an organization can either be a GitHub organization or a GitHub user
CREATE TABLE organizations (
    name VARCHAR(255) NOT NULL PRIMARY KEY,
    plan VARCHAR(255) NOT NULL,
    members JSON,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO organizations(name, plan) VALUES
('runway', 'enterprise'),
('monalisa', 'free'),
('lisamona', 'team');

# the commands table
CREATE TABLE commands (
    id VARCHAR(255) NOT NULL PRIMARY KEY,
    organization VARCHAR(255) NOT NULL,
    repository VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    data JSON,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO commands(id, organization, repository, name, data) VALUES
('5ecfdb3a-c229-4982-b5b0-5cc87b8a616a', 'runwayapp', 'test-flight', 'deploy command', '{"name": "deploy command", "state": "active", "description": "Deploy the application", "command": ".deploy", "actions": [{"type": "reaction", "mode": "add", "reaction": "+1"}]}'),
('8ff93daa-66dc-4398-9ad7-93a480ac8ad7', 'runwayapp', 'test-flight', 'linter', '{"name": "linter", "description": "it lints things", "command": ".lint", "state": "active", "actions": []}'),
('33393daa-66dc-4398-9ad7-93a480ac8333', 'runwayapp', 'test-flight', 'test command', '{"name": "test command", "description": "triggers an Actions workflow and leaves a comment", "command": ".test", "state": "active", "actions": [{"type": "comment", "text": "I am starting the [.github/workflows/test.yml](https://github.com/runwayapp/test-flight/actions/workflows/test.yml) workflow via a dispatch"}, {"type": "workflow_dispatch", "path": "test.yml"}]}'),
('58890287-9ff4-4ffa-b671-05ac33b9372e', 'runwayapp', 'fake-repo', 'help', '{"name": "help", "description": "a general help command", "command": ".help", "state": "active"}'),
('5a253c4d-ae3f-4b8d-aabd-f418c34f1d1f', 'monalisa', 'cats', 'help', '{"name": "help", "description": "a general help command", "command": ".help", "state": "active", "actions": []}'),
('e497b87c-7bc7-4565-8477-54c8f9441cd0', 'lisamona', 'dogs', 'help', '{"name": "help", "description": "a general help command", "command": ".help", "state": "active", "actions": []}');

# the users table
CREATE TABLE users (
    login VARCHAR(255) NOT NULL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO users(login) VALUES
('maverick'),
('goose');

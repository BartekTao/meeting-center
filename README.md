# meeting-center-backend

## branch policy
- Do not merge into `main` without a PR and a review.
- When developing any feature, create a new branch named `feat/$(your_branch_name)`.

## set up
### required
- install docker
- install docker-compose

### start backend service
#### part 1: run server
- `cd backend`
- copy tmp.env to .env, then set `SENDER_EMAIL`, `SENDER_PASSWORD`, `GOOGLE_OAUTH_CLIENT_SECRET` and `GOOGLE_OAUTH_CLIENT_ID` on .env
- download GCS credential json file, put it under `backend` folder (IMPORTANT!!)
- `docker build -f .\deployment\Dockerfile -t meeting-center-api .`
- `docker-compose -f ./deployment/Docker-compose.yml up -d`
- remove cli: `docker-compose -f ./deployment/Docker-compose.yml down`
#### part 2: use login
- call `http://localhost:8080/auth/google/login` to get token
- open `http://localhost:8080` to see graphQL doc
- set `{ "Authorization": "Bearer your-token" }` on header to call other api

#### api test
- open `http://localhost:8080`
- special case: uploadFile:
    - `echo "This is a test file" > testfile.txt`
    - 
        ```
        curl -X POST http://localhost:8080/query \
        -H "Authorization: Bearer $(your-token)" \
        -F operations='{ "query": "mutation ($file: Upload!) { uploadFile(file: $file) }", "variables": { "file": null } }' \
        -F map='{ "0": ["variables.file"] }' \
        -F 0=@testfile.txt
        ```
    - see more on https://gqlgen.com/reference/file-upload/

### start frontend service

#### Project setup
```
cd frontend
npm install
```

#### Run the serve
```
npm run serve
```

Visit http://localhost:8888

### Enter the service
Enter http://localhost:8888/#/login-page to get the token by a Google account.
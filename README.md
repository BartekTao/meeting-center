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
- copy tmp.env to .env, then set GOOGLE_OAUTH_CLIENT_SECRET and GOOGLE_OAUTH_CLIENT_ID on .env
- `docker build -f .\deployment\Dockerfile -t meeting-center-api .`
- `docker-compose -f ./deployment/Docker-compose.yml up -d`
- remove cli: `docker-compose -f ./deployment/Docker-compose.yml down`
#### part 2: use login
- call `http://localhost:8080/auth/google/login` to get token
- open `http://localhost:8080` to see graphQL doc
- set `{ "Authorization": "Bearer your-token" }` on header to call other api

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
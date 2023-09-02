# avito-segmentation
# WORK IN PROGRESS (MINIMUM DONE)

## Build
#### Prerequisites
- docker

1. Clone project:
```
git clone https://github.com/vlasashk/avito-segmentation.git
cd avito-segmentation
```
2. Run:
```
docker compose up --build
```
## Project information
API for dynamic user segmentation for testing new functionality
### Storage
- PostgreSQL
### Functionality
#### Users manipulation
- {POST} **/user/new** - Add new user to database.</br> Request Body JSON:
```
{
    "user_id": 10
}
```
- {POST} **/user/addSegment** - Add list of segments to user. 
User and each segment must be present in database for successful execution 
otherwise it won't ve allowed.</br> Request Body JSON:
```
{
    "user_id": 10,
    "segment_slug": ["AVITO","AVITO_10", "AVITO_30"]
}
```
- {GET} **/user/segments** - Return the list of segments the user is a member of.</br> Request Body JSON:
```
{
    "user_id": 10
}
```
- {DELETE} **/user/segments** - Remove user from chosen segments by marking deleted_at field.</br> Request Body JSON:
```
{
    "user_id": 10,
    "segment_slug": ["AVITO","AVITO_10", "AVITO_30"]
}
```

#### Segments manipulation
- {POST} **/segment/new** Add new segment to database.</br> Request Body JSON:
```
{
    "slug": "test"
}
```
- {DELETE} **/segment/remove**  Cascade delete segment. 
This method will permanently delete segment and all it's relations between user-segment.</br> Request Body JSON:
```
{
    "slug": "test"
}
```
- {GET} **/segment/users** - Return the list of users the segment has.</br> Request Body JSON:
```
{
    "slug": test
}
```
#### CSV Report
- {GET} **/report** - Return the link to csv file with report for chosen month (currently returns local path).</br> Request Body JSON:
```
{
    "year": 2023,
    "month": 9
}
```
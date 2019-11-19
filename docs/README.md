# Hemp Conduit API
wip

## Domain Functions
- user
    - [ ] register
    - [X] login
    - [X] update
- product  
    - [X] save
    - [X] find-all
    - [X] find-by-owner
- address
    - [X] save
    - [X] find-by-ids
- order
    - [ ] save
    - [ ] find-by-ids


## References
+ **go(lang)** - [simple, reliable, efficient][^go]
+ **aws** 
    + cli - [amazon web services (cli)][^aws]
    + sam - [server-less application model][^sam]
    + λƒ - [aws lambda function][^λƒ]
    + go sdk - [aws sdk for go api reference][^sdk]
    + DynamoDB - [NoSQL database service][^ddb]
+ **yq** - [command-line YAML processor][^yq]
+ **jq** - [command-line JSON processor][^jq]
+ **jwt** - [JSON web tokens][^jwt]

***

[^sam]: https://github.com/awslabs/serverless-application-model
[^sdk]: https://docs.aws.amazon.com/sdk-for-go/api/aws/
[^λƒ]: https://docs.aws.amazon.com/cli/latest/reference/lambda/index.html
[^go]: https://golang.org/
[^yq]: http://mikefarah.github.io/yq/
[^jq]: https://stedolan.github.io/jq/
[^aws]: https://aws.amazon.com/cli/
[^jwt]: https://jwt.io/
[^ddb]: https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Introduction.html
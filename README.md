
# GO SAM API

Source code contained herein provides a framework for developing and maintaining a tiered network of loosely coupled, 
highly cohesive, server-less micro-services. 

## Layered Architecture

A service oriented architecture is dominant given that service requests and responses of a server-less application model 
are materialized through message oriented middleware.

![Image](assets/ddd.jpg?raw=true)

## Services

Conceptually, (Micro)services are organized around a single branch of the business domain model and provide discrete 
units of functionality, to achieve predefined business objectives by working alone or with sibling services.

Practically, these are characterized as fine grained and independently deployable, capable of asynchronously 
facilitating decentralized HTTP requests synonymous with modern eCommerce platforms.

## Entities, Values, Aggregates

## License

GO SAM API is [MIT licensed](./LICENSE). By contributing to GO SAM API, you agree that your contributions will be 
licensed under its MIT license.

## todo - inline urls vvv
+ **go(lang)** - [simple, reliable, efficient][^go]
+ **aws** 
  + cli - [amazon web services][^aws]
  + sam - [server-less application model][^sam]
  + λƒ - [aws lambda function][^λƒ]
  + go sdk - [aws sdk for go api reference][^sdk]
  + DynamoDB - [NoSQL database service][^ddb]
+ **jq** - [command-line JSON processor][^jq]
+ **jwt** - [JSON web tokens][^jwt]
+ **ups** - [UPS Developer Kit][^ups]
+ **ups** - [UPS Developer Kit][^ups]
+ **ups** - [UPS Developer Kit][^ups]
+ **api** - [Application Programming Interface][^api]

***

[^api]: https://www.google.com/search?q=api
[^sam]: https://github.com/awslabs/serverless-application-model
[^sdk]: https://docs.aws.amazon.com/sdk-for-go/api/aws/
[^λƒ]: https://docs.aws.amazon.com/cli/latest/reference/lambda/index.html
[^go]: https://golang.org/
[^yq]: http://mikefarah.github.io/yq/
[^jq]: https://stedolan.github.io/jq/
[^aws]: https://aws.amazon.com/cli/
[^jwt]: https://jwt.io/
[^ddb]: https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Introduction.html
[^ups]: https://www.ups.com/upsdeveloperkit/announcements
[^soa]: https://en.wikipedia.org/wiki/Service-oriented_architecture
[^mic]: https://en.wikipedia.org/wiki/Microservices
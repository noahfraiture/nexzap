let payload = {
  LanguageName: "Go",
  TestContents: ["You must write a good test", "You must write a correct test"],
  GuideContents: [(cat 1_overview.md), (cat 2_struct.md)],
  DockerImages: ["gotest", "gotest"]
}

http post http://localhost:8080/api/tutorials $payload --content-type application/json

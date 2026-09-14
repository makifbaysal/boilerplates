# ios-native boilerplate

SwiftUI + Swift 6 Task CRUD app calling the same REST shape every
`backend/*` boilerplate exposes.

Read [.ai/architecture.md](.ai/architecture.md) first — it documents
three real bugs hit building this (most importantly: the domain model
is named `TaskItem`, not `Task`, because `Task` collides with Swift's
own concurrency primitive) so they don't regress.

## Commands

```
xcodegen generate            # regenerate TaskTrooperIOS.xcodeproj (gitignored)
                                # re-run after adding/removing/renaming any source file

xcodebuild build -project TaskTrooperIOS.xcodeproj -scheme TaskTrooperIOS \
  -destination 'platform=iOS Simulator,name=iPhone 17'

xcodebuild test -project TaskTrooperIOS.xcodeproj -scheme TaskTrooperIOS \
  -destination 'platform=iOS Simulator,name=iPhone 17'
```

`AppConfig.defaultAPIBaseURL` defaults to `http://localhost:8080` — see
architecture.md#config.

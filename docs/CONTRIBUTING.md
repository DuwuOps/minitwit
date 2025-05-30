# Contribution Guidelines

This document describes the agreed upon contribution guidelines by group "*DuwuOps*"' regarding their MiniTwit DevOps project.

## Tasks & Cooperation
* **Note:** Was changed in later iteration by unanimous vote, but served as guiding principle until halfway through development - as such, one contributor may complete a task by themselves.
1. A task must be completed through *cooperation*.
2. A task must be assigned to *at least* 2 contributors, who share the responsibility of it's completion.
3. The process of completing a task should, unless specified otherwise, be completed through [*pair-programming*](https://qentelli.com/thought-leadership/insights/introduction-pair-programming#:~:text=Pair%20programming%20is%20a%20concept,the%20accuracy%20of%20the%20code). 
    * Unless it is deemed unecessary for a task, in which case, a task will be assigned the label *exempt from pair programming*.
    

## GitHub Issues

1. Tasks must be represented by an appropriate GitHub issue.
2. A GitHub issue must be moved to the column of *‘done’* on completion of it's contents and the merging of it's associated branch.
3. A GitHub issue must include a list of *acceptance-criteria*.
4. A GitHub issue must be assigned an appropriate *label* upon creation.
5. A GitHub issue must be associated with a *time-spent* notion upon completion.

## Branches

* All branches must be associated to a **task** (GitHub issue).
* Pull-Requests must be approved by at least **2** other group members before a merge can commence.
    * These **2** members should not be part of the task's assigned contributors. 



# <a name="commit"></a> Git Commit Guidelines

The work of a task must be split into several commits of an appropriate size.

## <a name="commit-message-format"></a> Commit Message Format

Each commit message consists of a **header**, a **body**, and (potentially) **co-authors**. The header has a special
format that includes a **type** and a potential **scope**:

```html
    <type> (<scope>) : <body>

    <co-authors>
```

### Type

Must be one of the following:

* **Feat** : A new feature
* **Fix** : A bug fix
* **Docs** : Documentation only changes
* **Style** : Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc)
* **Refactor** : A code change that neither fixes a bug nor adds a feature
* **Perf** : A code change that improves performance
* **Test** : Adding missing tests
* **Chore** : Changes to the auxiliary tools such as release scripts
* **Build** : Changes to the dependencies, devDependencies, or build tooling
* **Ci** : Changes to our Continuous Integration configuration

### Scope

The scope could be anything that helps specify the scope (or feature) that is changing.

Examples
- Fix (select) : 
- Docs (menu): 

### Subject



### Body

The body starts with a **subject**. The subject contains a succinct description of the change.

* use the imperative, present tense: "change" not "changed" nor "changes"
* Capitalize first letter


Just as in the **subject**, use the imperative, present tense: "change" not "changed" nor "changes".
The body should include the motivation (the why) for the change and contrast this with previous behavior.

The body can also include a **footer**.
The footer should contain any information about **Breaking Changes** and is also the place to
reference GitHub issues that this commit **Closes**, **Fixes**, or **Relates to**.

> We highlight Breaking Changes in the ChangeLog. These are as changes that will require
  community users to modify their code after updating to a version that contains this commit.


### Co-authors

You can attribute a commit to more than one author by adding one or more Co-authored-by trailers to the commit's message.
Co-authors are written as the following:
`Co-authored-by: name <name@example.com>`

#### LLMs

Partly due to the [requirements of the course](https://github.com/itu-devops/lecture_notes/blob/master/sessions/session_02/README_TASKS.md#guidelines-for-using-ai-code-assistants-like-github-copilot-chatgpt-etc), LLMs are registered as co-authors as soon as you reused the tiniest bit of what they tell you while you are working on a task.

Note that the email is not necessary for LLMs

For example:
```text
    Co-authored-by: ChatGPT <openai@github.com>
    Co-authored-by: DeepSeek <deepseek-ai@github.com>
    Co-authored-by: GitHub Copilot <githubcopilot@github.com>
```



---

### Sample Commit messages

For example:
```text
    Fix (frontend-backend-connection) : The frontend pointed to the wrong IP for the backend.

    Co-authored-by: NAME <NAME@EXAMPLE.COM>
    Co-authored-by: ANOTHER-NAME <ANOTHER-NAME@EXAMPLE.COM>
```

```text
    Fix (autocomplete) : Don't show the menu panel when readonly

    - this could sometimes happen when no value was selected

    Fixes #11231
```

```text
    Feat (chips) : Trigger ng-change on chip addition/removal

    - add test of `ng-change` for `md-chips`
    - add docs regarding `ng-change` for `md-chips` and `md-contact-chips`
    - add demo for ng-change on `md-chips`
    - add demo for ng-change on `md-contact-chips`

    Fixes #11161 Fixes #3857
```

```text
    Refactor (content) : Prefix mdContent scroll- attributes

    BREAKING CHANGE: md-content's `scroll-` attributes are now prefixed with `md-`.

    Change your code from this:
    ```html
    <md-content scroll-x scroll-y scroll-xy>
    ```

    To this:
    ```html
    <md-content md-scroll-x md-scroll-y md-scroll-xy>
    ```
```
<br/>

## Coding Conventions

* Formatting Rules - *TBD*


## Acknowledgements
The section about git commits is based on [angular/.github/CONTRIBUTING.md](https://github.com/angular/material/blob/master/.github/CONTRIBUTING.md), but has been modified for clarity and relevance to this project.
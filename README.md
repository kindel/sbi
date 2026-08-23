# sbi

An interactive app that structures feedback as Situation, Behavior, and Impact so you leave with words about what happened, not who they are.

## The Model

SBI keeps feedback factual and focused. Three parts:

1. **Situation**: When and where. Specific enough to jog memory without debate. "During our team meeting on Tuesday at 10 a.m." beats "last week."
2. **Behavior**: What you saw or heard. Actions, not intent. "You raised your voice and interrupted the speaker," not "You were rude."
3. **Impact**: The effect on you, others, or the work.

The Center for Creative Leadership (CCL) is where the model originated. CCL adds a fourth beat: **Intent**. After the impact, ask "What were you hoping to accomplish?" That turns a statement into a conversation.

## Why it works

Most people already know they are not perfect and have behaviors they can change. SBI isolates behavior from character. Feedback about what someone did, not who they are, is harder to hear as an attack.

It works for praise as well as correction. Positive SBI examples become the evidence in a promotion case. A trail of SBI notes replaces gut feel in a review.

## The App

A single page that turns a vent into something you can say out loud. Paste the draft you'd actually deliver; the app highlights judgment words, mind-reading, absolutes, and vague timing and explains why each one weakens the feedback. Then rebuild it in three guided fields (Situation, Behavior, Impact), assemble the script, and end with the intent question. Copy it and go say it.

Nothing you type leaves the page. No server-side logic, no storage, no account.

### Running locally

```sh
go run .
```

Open http://localhost:8080. The app is the static files under `static/`, embedded into the binary; asset paths are relative so it works at the root locally and under a path prefix in production.

### Canonical path

The public URL is https://kindel.com/apps/sbi/. Legacy aliases `/sbi/` and `/tools/sbi/` also work.

## Teaching

- [The Secret to Giving Feedback Without Triggering Defensiveness](https://blog.kindel.com/2025/10/21/the-secret-to-giving-feedback-without-triggering-defensiveness/)

## Related

- [Office Hours](https://kindel.com/officehours/)

## License

MIT. Copyright (c) 2026 Kindel, LLC. Keep the copyright notice and permission notice in all copies.

Derivatives must include a visible link to https://kindel.com.

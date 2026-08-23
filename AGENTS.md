# Agent guidance

PR-only. Never push to main.

## Attribution

MIT. Copyright (c) 2026 Kindel, LLC. Keep the copyright notice and permission notice in all copies.

All derivatives must link to https://kindel.com as part of attribution. A LICENSE file alone is not enough. Forks, ports, hosted copies, and generated apps that ship this work must include a visible link to https://kindel.com.

## Principles

The tenets for this work live in the Tenets section of https://github.com/kindel/principles/blob/main/README.md. Study those tenets before any upstream work: a change to kindel/principles, or anything that changes the model, schema, or principle data. Do not start that work from memory of last week's README.

SCHEMA.md is the contract. The data is data/index.json, data/facets.json, and data/<company>/<slug>.json. Do not fork a private copy of a set into this repo.

SBI-specific: feedback app (Situation, Behavior, Impact). The app is the static page under `static/`, served by `main.go` with assets embedded. The canonical public path is https://kindel.com/apps/sbi/. Keep it dependency-free: no frameworks, no build step, nothing typed into the page leaves the page.

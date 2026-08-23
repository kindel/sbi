"use strict";

// Each rule flags language that makes feedback about the person instead of
// what happened. Phrases are matched case-insensitively on word boundaries.
const RULES = [
  {
    cat: "label",
    msg: "A word about who they are, not what they did. Swap it for what you saw or heard.",
    phrases: [
      "rude", "lazy", "unprofessional", "disrespectful", "arrogant",
      "careless", "sloppy", "incompetent", "toxic", "difficult",
      "aggressive", "condescending", "immature", "selfish", "unreliable",
      "annoying", "mean", "hostile", "dramatic", "terrible", "awful",
      "horrible", "a jerk", "bad attitude", "attitude",
    ],
  },
  {
    cat: "praise",
    msg: "Kind, but not evidence. What exactly did they do, and what did it make possible?",
    phrases: [
      "great", "amazing", "awesome", "fantastic", "excellent", "brilliant",
      "incredible", "outstanding", "good job", "nice work", "a rockstar",
    ],
  },
  {
    cat: "intent",
    msg: "That's a guess about what's in their head. Describe the action; save intent for the question at the end.",
    phrases: [
      "on purpose", "deliberately", "intentionally", "don't care",
      "doesn't care", "didn't care", "you were trying", "you're trying",
      "you think", "you thought", "you meant", "you wanted",
      "you just want", "your intention", "clearly you", "obviously you",
    ],
  },
  {
    cat: "absolute",
    msg: "\u201CAlways\u201D and \u201Cnever\u201D invite a debate about the one exception. Name one specific time.",
    phrases: ["always", "never", "constantly", "every time", "all the time"],
  },
  {
    cat: "vague",
    msg: "Specific enough to jog memory without debate \u2014 a date, a meeting, a document.",
    phrases: [
      "recently", "last week", "the other day", "a while ago",
      "a while back", "lately", "sometimes", "often", "in general",
      "a few times", "as of late",
    ],
  },
];

for (const rule of RULES) {
  const escaped = rule.phrases
    .slice()
    .sort((a, b) => b.length - a.length)
    .map((p) => p.replace(/[.*+?^${}()|[\]\\]/g, "\\$&").replace(/ /g, "\\s+"));
  rule.re = new RegExp("\\b(?:" + escaped.join("|") + ")\\b", "gi");
}

function lint(text) {
  const found = [];
  for (const rule of RULES) {
    rule.re.lastIndex = 0;
    let m;
    while ((m = rule.re.exec(text)) !== null) {
      found.push({ start: m.index, end: m.index + m[0].length, text: m[0], rule });
    }
  }
  // Keep the earliest match where rules overlap (e.g. "bad attitude" vs "attitude").
  found.sort((a, b) => a.start - b.start || b.end - a.end);
  const flags = [];
  let last = 0;
  for (const f of found) {
    if (f.start >= last) {
      flags.push(f);
      last = f.end;
    }
  }
  return flags;
}

function escapeHTML(s) {
  return s.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
}

function renderBackdrop(backdrop, text, flags) {
  let html = "";
  let pos = 0;
  for (const f of flags) {
    html += escapeHTML(text.slice(pos, f.start));
    html += '<mark class="' + f.rule.cat + '">' + escapeHTML(f.text) + "</mark>";
    pos = f.end;
  }
  // Trailing newline needs a visible stand-in or the backdrop ends one line short.
  html += escapeHTML(text.slice(pos)) + "\n";
  backdrop.innerHTML = html;
}

function renderFlags(list, text, flags) {
  list.textContent = "";
  const seen = new Set();
  for (const f of flags) {
    const key = f.text.toLowerCase() + "|" + f.rule.cat;
    if (seen.has(key)) continue;
    seen.add(key);
    const li = document.createElement("li");
    const chip = document.createElement("span");
    chip.className = "chip " + f.rule.cat;
    chip.textContent = f.text;
    li.appendChild(chip);
    li.appendChild(document.createTextNode(f.rule.msg));
    list.appendChild(li);
  }
  if (flags.length === 0 && text.trim() !== "") {
    const li = document.createElement("li");
    li.className = "ok";
    li.textContent = "\u2713 No judgment words \u2014 this is about what happened.";
    list.appendChild(li);
  }
}

function wireField(id) {
  const textarea = document.getElementById(id);
  const backdrop = textarea.closest(".lintbox").querySelector(".backdrop");
  const list = document.getElementById("flags-" + id);
  const update = () => {
    const text = textarea.value;
    const flags = lint(text);
    renderBackdrop(backdrop, text, flags);
    renderFlags(list, text, flags);
    backdrop.scrollTop = textarea.scrollTop;
    if (id !== "draft") renderScript();
  };
  textarea.addEventListener("input", update);
  textarea.addEventListener("scroll", () => {
    backdrop.scrollTop = textarea.scrollTop;
  });
  update();
}

const INTENT_QUESTION = "What were you hoping to accomplish?";

function buildScript() {
  const parts = ["situation", "behavior", "impact"]
    .map((id) => document.getElementById(id).value.trim())
    .filter((s) => s !== "");
  if (parts.length === 0) return "";
  let out = parts.join(" ");
  if (document.getElementById("intent").checked) {
    out += "\n\n" + INTENT_QUESTION;
  }
  return out;
}

function renderScript() {
  const block = document.getElementById("script");
  const copy = document.getElementById("copy");
  const text = buildScript();
  if (text === "") {
    block.innerHTML =
      '<span class="placeholder">Your feedback appears here as you fill in the three parts.</span>';
    copy.disabled = true;
  } else {
    block.textContent = text;
    copy.disabled = false;
  }
}

document.getElementById("intent").addEventListener("change", renderScript);

document.getElementById("copy").addEventListener("click", async () => {
  const text = buildScript();
  const copied = document.getElementById("copied");
  try {
    await navigator.clipboard.writeText(text);
    copied.textContent = "Copied. Say it soon, in private.";
  } catch {
    copied.textContent = "Couldn't copy \u2014 select the text above instead.";
  }
  setTimeout(() => {
    copied.textContent = "";
  }, 4000);
});

for (const id of ["draft", "situation", "behavior", "impact"]) {
  wireField(id);
}
renderScript();

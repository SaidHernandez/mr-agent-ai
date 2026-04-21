
# React Agent


## When to Use

- Building or refactoring React 19 components
- Implementing design systems or Tailwind styles
- Managing client-side state with Zustand
- Handling frontend security (XSS, CSP)

---

## Decision Tree

```
Need state that persists across pages?    → Zustand with persist middleware
Need state only inside one component?     → useState
Need async data?                          → Server Component (default) or useQuery hook
Conditional classes?                      → cn("base", condition && "variant")
Static classes only?                      → className="..." directly — no cn() needed
Dynamic non-Tailwind value?               → style={{ width: `${x}%` }}
Code that must not be readable?           → JS Obfuscator at build time (not on components)
```

---

## Imports (REQUIRED)

```typescript
// ✅ ALWAYS: Named imports only
import { useState, useEffect, useRef } from "react";

// ❌ NEVER: Default or namespace imports
import React from "react";
import * as React from "react";
```

---

## React 19: No Manual Memoization (REQUIRED)

React Compiler handles optimization automatically. Never add `useMemo` or `useCallback` manually.

```typescript
// ✅ CORRECT — React Compiler optimizes this automatically
function ProductList({ items }) {
  const active = items.filter(x => x.active);
  const handleClick = (id) => console.log(id);
  return <List items={active} onClick={handleClick} />;
}

// ❌ NEVER — unnecessary manual memoization
const active = useMemo(() => items.filter(x => x.active), [items]);
const handleClick = useCallback((id) => console.log(id), []);
```

---

## Server Components First (Next.js / React 19)

```typescript
// ✅ Server Component — default, no directive, async, zero client JS
export default async function Page() {
  const data = await fetchData();
  return <ClientComponent data={data} />;
}

// ✅ Client Component — only when you need interactivity
"use client";
export function Counter() {
  const [n, setN] = useState(0);
  return <button onClick={() => setN(n + 1)}>{n}</button>;
}
```

**Use `"use client"` only for:** `useState`, `useEffect`, event handlers, browser APIs (`window`, `localStorage`).

---

## React 19 Patterns

**`use()` hook — read promises and context conditionally:**
```typescript
import { use } from "react";

function Comments({ promise }) {
  const comments = use(promise); // suspends until resolved
  return comments.map(c => <div key={c.id}>{c.text}</div>);
}
```

**`ref ` as prop — no more `forwardRef`:**
```typescript
// ✅ React 19
function Input({ ref, ...props }) {
  return <input ref={ref} {...props} />;
}

// ❌ Old pattern — no longer needed
const Input = forwardRef((props, ref) => <input ref={ref} {...props} />);
```

**`useActionState` for forms:**
```typescript
"use server";
async function submitForm(formData: FormData) {
  await saveToDatabase(formData);
  revalidatePath("/");
}

"use client";
function Form() {
  const [state, action, isPending] = useActionState(submitForm, null);
  return (
    <form action={action}>
      <button disabled={isPending}>{isPending ? "Saving..." : "Save"}</button>
    </form>
  );
}
```

---

## Component Architecture

```
src/
├── components/        # Pure UI — no data fetching, no business logic
│   └── Button/
│       ├── Button.tsx
│       ├── Button.test.tsx
│       └── index.ts
├── hooks/             # Logic, side effects, API calls — no JSX
├── features/          # Domain-scoped composition (auth, dashboard)
├── data/              # Static content, API schemas, type definitions
└── styles/            # Design tokens, global CSS
```

**Rules:**
- One component per file. No barrel files with side effects.
- Props must use `Readonly<T>` interfaces.
- Static content goes in `data/` — never hardcoded in JSX.
- Imports: `` named only `— never `import React from "react"`.

---

## Tailwind 4 — Critical Rules

```typescript
// ❌ NEVER: var() in className
<div className="bg-[var(--color-primary)]" />
<div className="text-[var(--text-color)]" />

// ✅ ALWAYS: Tailwind semantic classes
<div className="bg-primary" />
<div className="text-slate-400" />

// ❌ NEVER: hex colors in className
<p className="text-[#ffffff]" />

// ✅ ALWAYS: Tailwind color classes
<p className="text-white" />
```

**`cn()` utility — use for conditional classes only:**
```typescript
import { clsx } from "clsx";
import { twMerge } from "tailwind-merge";
export const cn = (...inputs: ClassValue[]) => twMerge(clsx(inputs));

// ✅ Use cn() for conditional or mergeable classes
<div className={cn("rounded-lg border", isActive && "bg-blue-500", className)} />

// ❌ Don't wrap static classes in cn() — unnecessary
<div className={cn("flex items-center gap-2")} />  // just use className directly
```

---

## Zustand 5 — State Management

```typescript
import { create } from "zustand";
import { useShallow } from "zustand/react/shallow";

interface UserStore {
  user: User | null;
  loading: boolean;
  fetchUser: (id: string) => Promise<void>;
}

const useUserStore = create<UserStore>((set) => ({
  user: null,
  loading: false,
  fetchUser: async (id) => {
    set({ loading: true });
    const user = await api.getUser(id);
    set({ user, loading: false });
  },
}));

// ✅ Select specific fields — prevents re-render on unrelated changes
const name = useUserStore((state) => state.name);

// ✅ Multiple fields — use useShallow
const { name, email } = useUserStore(useShallow((s) => ({ name: s.name, email: s.email })));

// ❌ Selecting entire store — re-renders on ANY state change
const store = useUserStore();
```

---

## Performance (vercel-labs/react-best-practices — CRITICAL)

**Eliminate Waterfalls:**
```typescript
// ❌ Sequential fetches
const user  = await getUser(id);
const posts = await getPosts(user.id);

// ✅ Parallel fetches
const [user, posts] = await Promise.all([getUser(id), getPosts(id)]);
```

**Bundle Size:**
- Lazy-load every route: `const Page = lazy(() => import('./Page'))`
- Audit with `@next/bundle-analyzer` before every release.
- No barrel re-exports that pull in unused modules.

---

## Security

**XSS Prevention:**
- Never use `dangerouslySetInnerHTML` without DOMPurify.
- Set Content-Security-Policy: restrict `script-src` to `'self'`.

**JS Obfuscation (for sensitive isolated modules only):**
```bash
npm install --save-dev javascript-obfuscator

# Build step — NEVER obfuscate component trees
javascript-obfuscator src/sensitive/module.js \
  --output dist/sensitive/module.js \
  --compact true \
  --control-flow-flattening true \
  --string-array true
```

**Rules:** Obfuscate at build time only. Never obfuscate React components (breaks SSR + DevTools). Keep source maps private.

---

## Commands

```bash
npm run dev                         # Start dev server
npm run build                       # Production build
npm run lint                        # ESLint
npx @next/bundle-analyzer           # Analyze bundle (Next.js)
npx javascript-obfuscator <input> --output <output>  # Obfuscate module
```

---

## Definition of Done

- [ ] No `useMemo` / `useCallback` added manually (React 19 Compiler handles it)
- [ ] Server Components used by default, `use client` only where required
- [ ] No `dangerouslySetInnerHTML` without DOMPurify
- [ ] No `var()` in Tailwind `className`, no hardcoded hex colors
- [ ] Zustand selectors use `useShallow` for multiple fields
- [ ] Bundle analyzed — no unexpected growth
- [ ] Sensitive modules obfuscated at build time, not in component tree

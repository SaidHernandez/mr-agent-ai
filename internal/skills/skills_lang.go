package skills

// langSkills contains language and framework skills (single-select).
var langSkills = []Skill{
	{
		Name:        "react",
		Dir:         "lang/react",
		Description: "React 19, Zustand, Tailwind, Server Components.",
		Trigger:     "When building React 19 components, managing state with Zustand, or styling with Tailwind.",
		Content:     reactSkill,
	},
	{
		Name:        "vue",
		Dir:         "lang/vue",
		Description: "Vue 3, Composition API, Pinia, composables.",
		Trigger:     "When building Vue 3 components, managing state with Pinia, or structuring composables.",
		Content:     vueSkill,
	},
	{
		Name:        "typescript",
		Dir:         "lang/typescript",
		Description: "Strict mode, type inference, type guards.",
		Trigger:     "When writing TypeScript code, configuring tsconfig, or reviewing TS patterns.",
		Content:     typescriptSkill,
	},
	{
		Name:        "java",
		Dir:         "lang/java",
		Description: "Optional, null safety, clean API design.",
		Trigger:     "When writing Java code, designing APIs, or handling null/Optional patterns.",
		Content:     javaSkill,
	},
	{
		Name:        "python",
		Dir:         "lang/python",
		Description: "Type hints, Pythonic patterns, error handling.",
		Trigger:     "When writing Python code, designing modules, or reviewing Pythonic patterns.",
		Content:     pythonSkill,
	},
	{
		Name:        "golang",
		Dir:         "lang/golang",
		Description: "Error handling, goroutines, idiomatic Go.",
		Trigger:     "When writing Go code, designing packages, or reviewing idiomatic Go patterns.",
		Content:     golangSkill,
	},
}

// ─── Skill content functions ───────────────────────────────────────────────────

func reactSkill() string {
	return skillFrontmatter(
		"react",
		"React 19 agent. No manual memoization, Server Components first, Tailwind cn(), Zustand 5, named imports.",
		"When building React 19 components, managing state with Zustand, or styling with Tailwind.",
		"Read, Edit, Write, Glob, Grep, Bash",
	) + mustReadDoc("lang/react.md")
}

func vueSkill() string {
	return skillFrontmatter(
		"vue",
		"Vue 3 agent. Composition API with script setup, Pinia, composables, v-model, storeToRefs.",
		"When building Vue 3 components, managing state with Pinia, or structuring composables.",
		"Read, Edit, Write, Glob, Grep",
	) + mustReadDoc("lang/vue.md")
}

func typescriptSkill() string {
	return skillFrontmatter(
		"typescript",
		"TypeScript best practices. Strict mode, type inference, generics, async/await patterns.",
		"When writing TypeScript code, configuring tsconfig, or reviewing TS patterns.",
		"Read, Edit, Write, Glob, Grep",
	) + mustReadDoc("lang/typescript.md")
}

func javaSkill() string {
	return skillFrontmatter(
		"java",
		"Java best practices. Optional patterns, null safety, clean code, effective API design.",
		"When writing Java code, handling Optional, or designing safe APIs.",
		"Read, Edit, Write, Glob, Grep",
	) + mustReadDoc("lang/java.md")
}

func pythonSkill() string {
	return skillFrontmatter(
		"python",
		"Python best practices. Type hints, Pythonic patterns, error handling, performance tips.",
		"When writing Python code, designing modules, or reviewing Pythonic patterns.",
		"Read, Edit, Write, Glob, Grep, Bash",
	) + mustReadDoc("lang/python.md")
}

func golangSkill() string {
	return skillFrontmatter(
		"golang",
		"Go best practices. Error handling, interfaces, goroutines, idiomatic Go patterns.",
		"When writing Go code, designing packages, or reviewing idiomatic Go patterns.",
		"Read, Edit, Write, Glob, Grep, Bash",
	) + mustReadDoc("lang/golang.md")
}

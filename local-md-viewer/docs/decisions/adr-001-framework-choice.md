---
title: "ADR-001: Framework Choice"
---

# ADR-001: Framework Choice

## Status

**Accepted**

## Context

We need to choose a web framework for building the frontend application. The team has experience with multiple frameworks and we need to make a decision that balances developer productivity, performance, and long-term maintainability.

## Decision

We will use **Astro** as our primary framework for the following reasons:

### Pros

1. **Performance**: Zero JavaScript by default with optional hydration
2. **Flexibility**: Can use React, Vue, or Svelte components
3. **Developer Experience**: Excellent tooling and hot module replacement
4. **SEO**: Server-side rendering out of the box
5. **Content Focus**: Perfect for documentation and content-heavy sites

### Cons

1. Smaller ecosystem compared to Next.js
2. Team needs to learn new paradigms
3. Less mature for complex interactive applications

## Consequences

- Team will need 1-2 weeks of ramp-up time
- We can leverage existing React component knowledge
- Documentation site will load significantly faster
- Need to carefully consider when to add client-side JavaScript

## Alternatives Considered

1. **Next.js**: More mature but heavier for our use case
2. **Vite + React**: More familiar but no SSR without extra setup
3. **Hugo**: Fast but limited interactivity

## References

- [Astro Documentation](https://docs.astro.build)
- [Why Choose Astro](https://astro.build/blog/why-astro/)


# Java Best Practices


## When to Use

- Writing or reviewing Java code
- Designing service or repository layers
- Handling nullable return values
- API design decisions around null safety

**Don't use when:** Working in other JVM languages (Kotlin has its own null-safety model).

---

## Decision Tree

```text
Method might return no value?          → return Optional<T>
Value might be null from external src? → Optional.ofNullable()
Value guaranteed non-null?             → Optional.of()
Need a default for missing value?      → orElse() / orElseGet()
Missing value is an error?             → orElseThrow()
Want to transform the value?           → map() / flatMap()
Want to filter?                        → filter()
```

---

## Critical Patterns

### 1. Optional — Use for Return Types Only

```java
// ✅ GOOD — Optional as return type signals "might be absent"
public Optional<User> findById(String id) {
    return userRepository.findById(id);
}

// ❌ BAD — Optional as parameter adds unnecessary complexity
public void updateUser(Optional<String> name) { ... }

// ❌ BAD — Optional as field breaks JPA and serialization
public class User {
    private Optional<String> nickname; // never do this
}
```

### 2. Never Return null from Optional Methods

```java
// ✅ GOOD
public Optional<User> findByEmail(String email) {
    User user = db.query(email);
    return Optional.ofNullable(user); // safe even if user is null
}

// ❌ BAD — defeats the entire purpose of Optional
public Optional<User> findByEmail(String email) {
    return null; // throws NullPointerException on .get()
}
```

### 3. Avoid `isPresent()` + `get()` — Use Functional Style

```java
// ✅ GOOD — functional style
userRepo.findById(id).ifPresent(user -> send(user.getEmail()));

String name = userRepo.findById(id)
    .map(User::getName)
    .orElse("Anonymous");

// ❌ BAD — verbose imperative style
Optional<User> opt = userRepo.findById(id);
if (opt.isPresent()) {
    send(opt.get().getEmail());
}
```

### 4. `orElse()` vs `orElseGet()`

```java
// ✅ GOOD — orElseGet() is lazy: lambda runs only when absent
String value = optional.orElseGet(() -> expensiveComputation());

// ✅ GOOD — orElse() is fine for cheap/pre-computed defaults
String value = optional.orElse("default");

// ❌ BAD — orElse() with expensive call always evaluates it
String value = optional.orElse(expensiveComputation()); // called even when present!
```

### 5. `orElseThrow()` for Mandatory Values

```java
// ✅ GOOD — explicit, descriptive error
User user = userRepo.findById(id)
    .orElseThrow(() -> new UserNotFoundException("User not found: " + id));

// ❌ BAD — get() throws NoSuchElementException with no context
User user = userRepo.findById(id).get();
```

### 6. Chain `map()` and `flatMap()`

```java
// ✅ GOOD
Optional<String> city = userRepo.findById(id)
    .map(User::getAddress)
    .map(Address::getCity)
    .filter(c -> !c.isBlank());

// ❌ BAD — manual null checks with nesting
User user = userRepo.findById(id).orElse(null);
if (user != null && user.getAddress() != null) {
    String city = user.getAddress().getCity();
}
```

### 7. Avoid Nested Optionals

```java
// ✅ GOOD — flatMap flattens Optional<Optional<T>>
Optional<String> email = userRepo.findById(id)
    .flatMap(user -> contactRepo.findEmail(user.getId()));

// ❌ BAD — nested Optional is unergonomic
Optional<Optional<String>> nestedEmail = userRepo.findById(id)
    .map(user -> contactRepo.findEmail(user.getId()));
```

### 8. Don't Serialize Optionals in DTOs

```java
// ✅ GOOD — use nullable field in DTOs
public class UserDTO {
    private String nickname; // nullable, serializes cleanly to JSON
}

// ❌ BAD — Jackson doesn't serialize Optional well by default
public class UserDTO {
    private Optional<String> nickname; // serializes as {"present": true}
}
```

### 9. General Java Best Practices

```java
// ✅ Use final for immutable local references
final String userId = request.getId();

// ✅ Program to interfaces, not implementations
List<User> users = new ArrayList<>();  // declared as List, not ArrayList

// ✅ Use try-with-resources for AutoCloseable
try (var conn = dataSource.getConnection()) {
    // conn is auto-closed on exit
}

// ✅ Prefer records for immutable data carriers (Java 16+)
public record UserDTO(String id, String name, String email) {}

// ✅ Use var for local type inference where type is obvious (Java 10+)
var users = userRepository.findAll();
```

---

## Commands

```bash
./mvnw verify              # compile, test, and package
./mvnw spotbugs:check      # static analysis
./gradlew test             # run tests (Gradle)
```

---

## Definition of Done

- [ ] Optional used only as return type (not parameters or fields)
- [ ] No `get()` calls without prior `isPresent()` — prefer `orElseThrow()`
- [ ] `orElseGet()` used for expensive default computations
- [ ] No nested `Optional<Optional<T>>` — use `flatMap()`
- [ ] DTOs use nullable fields, not Optional fields
- [ ] No raw `null` returns from methods that should return Optional

---


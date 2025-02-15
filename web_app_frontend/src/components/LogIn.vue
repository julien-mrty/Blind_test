<template>
  <div class="connection-page">
    <h1>Login</h1>
    <form @submit.prevent="handleSubmit">
      <label for="username">Username:</label>
      <input v-model="username" id="username" type="text" required />

      <label for="password">Password:</label>
      <input v-model="password" id="password" type="password" required />

      <button type="submit">Login</button>
      <p v-if="errorMessage" class="error-message">{{ errorMessage }}</p>
    </form>
    <button @click="goToSignUp" class="signup-button">Sign Up</button>
  </div>
</template>


<script>
import { ref } from "vue"; // Import ref for reactive variables
import { useRouter } from "vue-router";

export default {
  setup() {
    const router = useRouter();

    // Reactive variables
    const username = ref("");
    const password = ref("");
    const errorMessage = ref("");

    const apiUrl = import.meta.env.VITE_APP_API_URL;

    const handleSubmit = async () => {
      try {
        const response = await fetch(`${apiUrl}/api/auth/login`, {

          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            username: username.value,
            password: password.value,
          }),
        });

        if (response.ok) {
          const data = await response.json();

          // Save the token, username and avatar_url to localStorage
          localStorage.setItem("authToken", data.token);
          localStorage.setItem("username", data.username);
          localStorage.setItem("avatar_url", data.avatar_url);

          console.log("Login successful:", data);
          console.log("Received Token:", data.token);


          // Redirect to the home page
          router.push("/home");
        } else {
          const error = await response.json();
          errorMessage.value = error.message || "Invalid username or password.";
        }
      } catch (error) {
        errorMessage.value = "An error occurred. Please try again.";
        console.error(error);
      }
    };

    const goToSignUp = () => {
      router.push("/signup");
    };

    return {
      username,
      password,
      errorMessage,
      handleSubmit,
      goToSignUp,
    };
  },
};
</script>


<style scoped>
.connection-page {
  max-width: 400px;
  margin: auto;
  padding: 20px;
  border: 1px solid #ccc;
  border-radius: 8px;
}

label {
  display: block;
  margin: 10px 0 5px;
}

input {
  width: 100%;
  padding: 8px;
  margin-bottom: 10px;
}

button {
  width: 100%;
  padding: 10px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

button[type="submit"] {
  background-color: #007bff;
  color: #fff;
}

button[type="submit"]:hover {
  background-color: #0056b3;
}

.signup-button {
  margin-top: 10px;
  background-color: #28a745;
  color: #fff;
}

.signup-button:hover {
  background-color: #218838;
}

.error-message {
  color: red;
  margin-top: 10px;
}
</style>
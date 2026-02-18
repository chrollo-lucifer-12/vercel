"use client";

import AuthForm from "@/components/custom/auth-form";
import { useSignIn } from "@/hooks/use-auth";
import { EnvelopeIcon, PasswordIcon } from "@phosphor-icons/react";

const SigninPage = () => {
  const { mutate, isPending, data } = useSignIn();
  return (
    <AuthForm
      title="Login to account"
      description="Login to your existing account"
      submitText="Login"
      isPending={isPending}
      loadingText="Logging in..."
      errors={data?.error}
      onSubmit={(formData) => mutate(formData)}
      footerText={
        <p>
          Don't have an account? <a href="/signup">Sign up</a>
        </p>
      }
      fields={[
        {
          id: "email",
          text: "Email",
          type: "email",
          placeholder: "john@gmail.com",
          icon: EnvelopeIcon,
          required: true,
        },
        {
          id: "password",
          text: "Password",
          type: "password",
          icon: PasswordIcon,
          required: true,
          placeholder: "",
        },
      ]}
    />
  );
};

export default SigninPage;

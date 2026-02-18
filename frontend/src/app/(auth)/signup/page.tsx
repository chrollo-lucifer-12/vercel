"use client";

import AuthForm from "@/components/custom/auth-form";
import { useSignUp } from "@/hooks/use-auth";
import {
  EnvelopeIcon,
  IdentificationCardIcon,
  PasswordIcon,
} from "@phosphor-icons/react";

const SignupPage = () => {
  const { mutate, isPending, data } = useSignUp();
  return (
    <AuthForm
      title="Create an account"
      description="Enter your information below to create your account"
      submitText="Create Account"
      loadingText="Creating account..."
      isPending={isPending}
      errors={data?.error}
      onSubmit={(formData) => mutate(formData)}
      footerText={
        <p>
          Already have an account? <a href="/signin">Sign in</a>
        </p>
      }
      fields={[
        {
          id: "name",
          text: "Full Name",
          type: "text",
          placeholder: "John Doe",
          icon: IdentificationCardIcon,
          required: true,
        },
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

export default SignupPage;

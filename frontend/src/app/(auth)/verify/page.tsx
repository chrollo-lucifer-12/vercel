"use client";

import AuthForm from "@/components/custom/auth-form";
import { useVerify } from "@/hooks/use-auth";
import { EnvelopeIcon } from "@phosphor-icons/react";

const VerifyPage = () => {
  const { mutate, isPending, data } = useVerify();
  return (
    <AuthForm
      title="Verify your account"
      description=""
      loadingText="Verifying..."
      submitText="Send Mail"
      isPending={false}
      errors={data?.error}
      onSubmit={(formData) => mutate(formData)}
      footerText={
        <p>
          Already have an account? <a href="/signin">Sign in</a>
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
      ]}
    />
  );
};

export default VerifyPage;

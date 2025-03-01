import { Routes, Route, Link, useNavigate } from "react-router-dom";
import AuthPage from "./Pages/AuthPage";
import Dashboard from "./Pages/Home";
import PasswordResetForm from "./Pages/PasswordResetForm";
import ForgotPasswordPage from "./Pages/ForgetPassword";

export default function App() {
  return (
    <div className="flex flex-col  min-h-screen bg-gray-100">
      <Routes>
        <Route path="/" element={<AuthPage />} />
        <Route path="/dashboard" element={<Dashboard />} />
        <Route path="/forget-password" element={<ForgotPasswordPage/>}/>
        <Route path="/complete-forget-password" element={<PasswordResetForm/>}/>
      </Routes>
    </div>
  );
}




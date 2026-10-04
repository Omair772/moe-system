# 🌐 Website Structure — Ministry of Education System

## 1. Purpose & Scope

This document defines the **React application structure** for the Ministry of Education integrated system, covering page division, navigation, routing, and file organization. It follows the **Clean Architecture** principles from pro-skills and the component-driven development patterns from delegate-skills.

## 2. Project Structure

```
src/
├─ assets/              │
│  ├─ images/          │
│  └─ fonts/           │
├─ components/          │
│  ├─ common/           │
│  │  ├─ Button.tsx      │
│  │  ├─ Input.tsx       │
│  │  ├─ Table.tsx       │
│  │  ├─ Card.tsx       │
│  │  └─ LoadingSpinner.tsx │
│  ├─ students/         │
│  │  ├─ List.tsx        │
│  │  ├─ Form.tsx        │
│  │  └─ Table.tsx       │
│  ├─ courses/          │
│  │  └─ similar structure│
│  ├─ teachers/         │
│  ├─ parents/          │
│  └─ shared/           │
│     ├─ auth/          │
│     │  ├─ LoginForm.tsx│
│     │  └─ useAuth.ts   │
│     ├─ api/           │
│     │  └─ client.ts    │
│     └─ toast/         │
│        └─ useToast.ts  │
├─ hooks/               │
│  useAuth.ts           │
│  useStudents.ts       │
│  useCourses.ts        │
│  useToast.ts          │
│  useTable.ts          │
├─ pages/               │
│  ├─ Login.tsx         │
│  ├─ Register.tsx      │
│  ├─ Dashboard.tsx     │
│  ├─ StudentsPage.tsx  │
│  ├─ CoursesPage.tsx   │
│  ├─ TeachersPage.tsx  │
│  ├─ ParentPortal.tsx  │
│  ├─ ReportsPage.tsx   │
│  └─ NotFound.tsx      │
├─ stores/              │
│  useStudentStore.ts   │
│  useCourseStore.ts    │
│  useAuthStore.ts      │
├─ styles/              │
│  ├─ tailwind.config.js│
│  ├─ global.css        │
│  └─ dark-mode.css     │
├─ types/               │
│  ├─ api.ts            │
│  ├─ student.ts        │
│  ├─ course.ts         │
│  └─ index.ts          │
└─ routes/              │
   └─ index.tsx         │
```

## 3. Route Configuration (React Router v6)

### 3.1 Main Routes

```tsx
// src/routes/index.tsx
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import Login from "@/pages/Login";
import Register from "@/pages/Register";
import Dashboard from "@/pages/Dashboard";
import StudentsPage from "@/pages/StudentsPage";
import CoursesPage from "@/pages/CoursesPage";
import TeachersPage from "@/pages/TeachersPage";
import ParentPortal from "@/pages/ParentPortal";
import ReportsPage from "@/pages/ReportsPage";
import NotFound from "@/pages/NotFound";
import StudentDetail from "@/pages/StudentDetail";
import CourseDetail from "@/pages/CourseDetail";

export const router = createBrowserRouter([
  // Public Routes
  {
    path: "/",
    element: <Login />,
    errorElement: <NotFound />,
  },
  {
    path: "/register",
    element: <Register />,
    errorElement: <NotFound />,
  },

  // Protected Routes (require authentication)
  {
    path: "/dashboard",
    element: <Dashboard />,
    errorElement: <NotFound />,
    loader: async () => {
      // Check auth state
      const { data } = await api.get("/auth/me");
      return data;
    },
    children: [
      // Student Routes
      {
        path: "students",
        element: <StudentsPage />,
        errorElement: <NotFound />,
        children: [
          {
            path: ":id",
            element: <StudentDetail />,
            errorElement: <NotFound />,
          },
        ],
      },

      // Course Routes
      {
        path: "courses",
        element: <CoursesPage />,
        errorElement: <NotFound />,
        children: [
          {
            path: ":code",
            element: <CourseDetail />,
            errorElement: <NotFound />,
          },
        ],
      },

      // Teacher Routes
      {
        path: "teachers",
        element: <TeachersPage />,
        errorElement: <NotFound />,
      },

      // Parent Routes
      {
        path: "parent",
        element: <ParentPortal />,
        errorElement: <NotFound />,
      },

      // Admin Routes
      {
        path: "reports",
        element: <ReportsPage />,
        errorElement: <NotFound />,
        // Only SUPER_ADMIN or ADMIN roles
        evaluateAccess: async () => {
          const { role } = await api.get("/auth/role");
          return role === "SUPER_ADMIN" || role === "ADMIN";
        },
      },
    ],
  },

  // Redirect to login if not authenticated
  {
    path: "*",
    element: <NotFound />,
  },
]);

export default function App() {
  return <RouterProvider router={router} />;
}
```

### 3.2 Route Protection & Role-Based Access

```tsx
// src/routes/protected.tsx
import { Navigate, useNavigate } from "react-router-dom";
import { useAuth } from "@/hooks/useAuth";

export function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const { isAuthenticated, role } = useAuth();
  const navigate = useNavigate();

  if (!isAuthenticated) {
    navigate("/login", { replace: true });
    return null;
  }

  // Role-based routing
  const roleRoutes: Record<string, string[]> = {
    STUDENT: ["dashboard/students", "dashboard/courses"],
    TEACHER: ["dashboard/teachers", "dashboard/reports"],
    PARENT: ["parent/portal", "parent/children"],
    ADMIN: ["dashboard", "dashboard/reports"],
    SUPER_ADMIN: ["dashboard", "dashboard/reports", "admin/users"],
  };

  const allowedRoutes = roleRoutes[role as keyof typeof roleRoutes];
  if (allowedRoutes && !allowedRoutes.includes(window.location.pathname)) {
    navigate("/dashboard", { replace: true });
    return null;
  }

  return children;
}
```

## 4. Page Components Structure

### 4.1 Authentication Pages

```
Login.tsx
├─ LoginForm (common component)
│  ├─ Email Input (Zod validation)
│  ├─ Password Input
│  ├─ Remember Me checkbox
│  └─ Submit Button
└─ Page Layout
   ├─ Title/Logo
   ├─ Form Card
   └─ Links: Register, Password Reset
```

### 4.2 Dashboard (Role-Specific)

```
Dashboard.tsx
├─ Header
│  ├─ User Avatar
│  ├─ User Name
│  ├─ Role Display
│  └─ Notification Bell
├─ Sidebar Navigation
│  ├─ Dashboard (home)
│  ├─ Students (if TEACHER/ADMIN/SUPER_ADMIN)
│  ├─ Courses (if TEACHER/ADMIN/SUPER_ADMIN)
│  ├─ Attendance (if TEACHER/ADMIN)
│  ├─ Grades (if TEACHER/ADMIN)
│  ├─ Reports (if ADMIN/SUPER_ADMIN)
│  └─ Settings
├─ Main Content Area
│  ├─ Quick Stats Card
│  ├─ Recent Activity
│  └─ Navigation Shortcuts
└─ Footer
   ├─ Version
   ├─ Help Link
   └─ Logout Button
```

### 4.3 Students Page

```
StudentsPage.tsx
├─ Toolbar
│  ├─ Add Student Button
│  ├─ Search Input
│  ├─ Filters (Grade Level, Status)
│  └─ Export Button
├─ Table
│  ├─ Checkbox Column (select multiple)
│  ├─ Name Column
│  ├─ Email Column
│  ├─ Enrolled Courses Column
│  ├─ Attendance % Column
│  └─ Actions Column (View, Edit, Enroll/Withdraw)
├─ Pagination (bottom)
│  ├─ Page size selector (10/25/50/100)
│  ├─ Page numbers
│  └─ Previous/Next buttons
└─ Empty State
   ├─ Illustration
   ├─ "No students found"
   └─ "Add your first student"
```

### 4.4 Course Detail Page

```
CourseDetail.tsx
├─ Breadcrumbs
│  ├─ Home
│  ├─ Courses
│  └─ [Course Code]
├─ Header
│  ├─ Course Code (e.g., "CS101")
│  ├─ Course Title
│  └─ Credits Display
├─ Sections
│  ├─ Course Information
│  │  ├─ Description
│  │  ├─ Prerequisites
│  │  ├─ Credit Hours
│  │  └─ Department
│  ├─ Enrollment Stats
│  │  ├─ Total Students
│  │  ├─ Capacity
│  │  ├─ Fill Percentage
│  │  └─ Waitlist Count
│  ├─ Instructor
│  │  ├─ Name
│  │  ├─ Department
│  │  ├─ Email
│  │  └─ Office Hours
│  └─ Schedule
│     ├─ Days (Mon, Wed, Fri)
│     ├─ Time (09:00-10:30)
│     ├─ Room Number
│     └─ Type (Lecture/Tutorial/Lab)
├─ Enrollment Form
│  ├─ Student Search
│  ├─ Enroll Button
│  └─ Waitlist Toggle
└─ Related Courses
   ├─ Prerequisites
   ├─ Concurrent Courses
   └─ Follow-up Courses
```

## 5. Navigation Menu Design

### 5.1 Responsive Navigation

```
Desktop (min-width: 1024px)
┌─────────────────────────────────────────────┐
│  ┌─────────────┐  ┌────────────────────┐ │
│  │  EduLogo     │  │  Search & Notifications│ │
│  └─────────────┘  └────────────────────┘ │
│  ┌─────────────────────────────────────┐ │
│  │  Main Navigation                      │ │
│  │  ├── Students                        │ │
│  │  ├── Courses                         │ │
│  │  ├── Teachers                        │ │
│  │  ├── Attendance                      │ │
│  │  ├── Grades                          │ │
│  │  ├── Reports (Admin only)            │ │
│  │  └─ Settings                         │ │
│  └─────────────────────────────────────┘ │
│  ┌─────────────────────────────────────┐ │
│  │  User Menu                           │ │
│  │  ├── Profile                         │ │
│  │  ├── Help                            │ │
│  │  └─ Logout                           │ │
│  └─────────────────────────────────────┘ │
└─────────────────────────────────────────────┘
```

```
Mobile (max-width: 768px)
┌─────────────────────┐
│  ☰ Menu              │
│ ─────────────────────│
│  Students            │
│  Courses             │
│  Teachers            │
│  Attendance          │
│  Grades              │
│  Reports (if admin)  │
│  Profile             │
│  Help                │
│  Logout              │
└─────────────────────┘
│  ← Swipe left for more options
└─────────────────────┘
```

### 5.2 Breadcrumbs Component

```tsx
// src/components/common/Breadcrumbs.tsx
import React from "react";
import { Link } from "react-router-dom";

interface BreadcrumbsProps {
  items: { label: string; path: string }[];
  currentLabel: string;
}

export function Breadcrumbs({ items, currentLabel }: BreadcrumbsProps) {
  return (
    <nav aria-label="breadcrumb">
      <ol className="flex items-center space-x-1 text-sm text-gray-600">
        {items.map((item, index) => (
          <li key={index} className="flex items-center">
            {index > 0 && <span className="mx-1 text-gray-400">/</span>}
            {index < items.length - 1 ? (
              <Link
                href={item.path}
                className="hover:text-primary transition-colors"
              >
                {item.label}
              </Link>
            ) : (
              <span className="font-medium text-gray-900">{item.label}</span>
            )}
          </li>
        ))}
        <li aria-current="page"
            className="text-gray-500 whitespace-nowrap font-medium"
        >
          {currentLabel}
        </li>
      </ol>
    </nav>
  );
}
```

## 6. State Management & Data Fetching

### 5.1 API Client (Axios with Interceptors)

```tsx
// src/api/client.ts
import axios, { AxiosInstance, AxiosError } from "axios";

const api: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_URL || "http://localhost:8080/api/v1",
  timeout: 10000,
  headers: {
    "Content-Type": "application/json",
    "Accept": "application/json",
  },
});

// Request interceptor - add auth token
api.interceptors.request.use((config) => {
  const token = localStorage.getItem("token");
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Response interceptor - handle auth errors
api.interceptors.response.use(
  (response) => response.data,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem("token");
      localStorage.removeItem("userRole");
      window.location.href = "/login";
    }
    return Promise.reject(error);
  }
);

export default api;
```

### 5.2 Custom Hooks Pattern

```tsx
// src/hooks/useStudents.ts
import { useState, useEffect } from "react";
import api from "@/api/client";
import { useNavigate } from "react-router-dom";

interface Student {
  id: string;
  name: string;
  email: string;
  enrolledCourses: string[];
  attendancePercentage: number;
}

interface UseStudentsReturn {
  students: Student[];
  loading: boolean;
  error: string | null;
  fetchStudents: (filters?: { gradeLevel?: string; search?: string }) => Promise<void>;
  deleteStudent: (id: string) => Promise<void>;
}

export function useStudents(): UseStudentsReturn {
  const [students, setStudents] = useState<Student[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();

  const fetchStudents = async (filters?: {
    gradeLevel?: string;
    search?: string;
  }) => {
    setLoading(true);
    setError(null);
    try {
      const { data } = await api.get("/students", {
        params: filters,
      });
      setStudents(data.data);
    } catch (err) {
      const axiosError = err as AxiosError;
      setError(axiosError.response?.data?.message || "Failed to fetch students");
    } finally {
      setLoading(false);
    }
  };

  const deleteStudent = async (id: string) => {
    if (!window.confirm("Are you sure you want to delete this student?")) {
      return;
    }
    try {
      await api.delete(`/students/${id}`);
      await fetchStudents();
      // Show success toast
      // toast.success("Student deleted successfully");
    } catch (err) {
      const axiosError = err as AxiosError;
      setError(axiosError.response?.data?.message || "Failed to delete student");
    }
  };

  useEffect(() => {
    fetchStudents();
  }, []);

  return {
    students,
    loading,
    error,
    fetchStudents,
    deleteStudent,
  };
}
```

## 6. SEO Specification for Main Interface

### 6.1 Metadata Configuration

```tsx
// src/pages/_document.tsx (Next.js data attribute)
import Document, { DocumentContext } from "next/document";

export default function MyDocument() {
  return (
    <Document>
      <html manifest="/manifest.json">
        <head>
          <meta charset="utf-8" />
          <meta
            name="viewport"
            content="width=device-width, initial-scale=1, viewport-fit=cover"
          />
          <title>
            {process.env.NEXT_PUBLIC_APP_TITLE || "Ministry of Education System"}
          </title>
          <meta
            name="description"
            content={
              process.env.NEXT_PUBLIC_APP_DESCRIPTION ||
              "Integrated ministry education management system"
            }
          />
          <link
            rel="canonical"
            href={process.env.NEXT_PUBLIC_APP_URL || "https://edu.example.com"}
          />

          {/* Open Graph / Social Media */}
          <meta
            property="og:title"
            content={process.env.NEXT_PUBLIC_APP_TITLE || "Ministry of Education"}
          />
          <meta
            property="og:description"
            content="Integrated ministry education management system"
          />
          <meta
            property="og:type"
            content="website"
          />
          <meta
            property="og:url"
            content={process.env.NEXT_PUBLIC_APP_URL || "https://edu.example.com"}
          />
          <meta
            property="og:image"
            content="/assets/og-image.png"
          />
          <meta
            property="og:locale"
            content="ar-AE" /* Arabic UAE context */}
          />

          {/* Twitter Cards */}
          <meta
            name="twitter:card"
            content="summary_large_image"
          />
          <meta
            name="twitter:title"
            content={process.env.NEXT_PUBLIC_APP_TITLE || "Ministry of Education"}
          />
          <meta
            name="twitter:description"
            content="Integrated ministry education management system"
          />
          <meta
            name="twitter:image"
            content="/assets/og-image.png"
          />

          {/* Favicons */}
          <link rel="icon" href="/favicon.ico" />
          <link
            rel="apple-touch-icon"
            href="/assets/apple-touch-icon.png"
          />
          <link
            rel="manifest"
            href="/manifest.json"
          />
        </head>
        <body>
          <div id="__next">{children}</div>
        </body>
      </html>
    </Document>
  );
}
```

### 6.2 Structured Data (JSON-LD)

```tsx
// src/components/SEO/Schema.org.tsx
import React from "react";

export function SchemaOrgEducationSite() {
  return (
    <script type="application/ld+json">
      {
        "@context": "https://schema.org",
        "@type": "EducationalOrganization",
        "name": "Ministry of Education",
        "alternateName": "MOE System",
        "description": "Integrated system for managing students, courses, attendance, and grades",
        "url": "https://edu.example.com",
        "logo": "https://edu.example.com/assets/logo.png",

        "address": {
          "@type": "PostalAddress",
          "addressLocality": "Capital City",
          "addressCountry": "Country Name"
        },

        "educationalLevel": "K-12",
        "founded": "2024",
        "memberOf": {
          "@type": "Organization",
          "name": "Government"
        },

        "knowsAbout": [
          "Student Management",
          "Course Registration",
          "Attendance Tracking",
          "Grade Management",
          "Report Generation"
        ],

        "sameAs": [
          "https://www.gov.edu.example.com",
          "https://twitter.com/edu_example"
        ]
      }
    </script>
  );
}
```

### 6.3 SEO Best Practices Checklist

| Element | Implementation | Status |
|---------|---------------|--------|
| **Title Tag** | Unique per page, < 60 chars | ✅ |
| **Meta Description** | 150-160 chars, keyword-rich | ✅ |
| **Header Hierarchy** | H1 per page, H2-H6 logical order | ✅ |
| **Alt Text** | Descriptive for all images | ✅ |
| **URL Structure** | `/:section/:id`, SEO-friendly slugs | ✅ |
| **Internal Linking** | Breadcrumbs, related content | ✅ |
| **Load Speed** | Lighthouse > 90, p95 < 200ms | ✅ |
| **Mobile-First** | Responsive, meta viewport set | ✅ |
| **Structured Data** | JSON-LD Organization schema | ✅ |
| **Robots.txt** | Allow core crawling, block admin routes | ✅ |
| **Sitemap.xml** | Auto-generated, submitted to search consoles | ✅ |

### 6.4 Performance Optimization

```
1. Code Splitting (React.lazy + Suspense)
   - Route-level code splitting
   - Component-level lazy loading
   - Dynamic imports for heavy components

2. Image Optimization
   - Next.js Image component (if using Next)
   - WebP/AVIF formats
   - Responsive srcSet
   - Lazy loading (loading="lazy")

3. Caching Strategy
   - SWR/react-query for data caching
   - Static generation (SSG) for read-only pages
   - Incremental Static Regeneration (ISR)

4. Bundle Optimization
   - Terser minification
   - Tree shaking
   - Analyze with webpack-bundle-analyzer
   - Target: < 200KB initial JS load
```

## 7. Accessibility (a11y) Standards

| Requirement | Implementation | WCAG Level |
|-------------|---------------|------------|
| **Keyboard Navigation** | All interactive elements reachable via Tab | A |
| **Screen Reader Support** | ARIA labels, roles, descriptive text | AA |
| **Color Contrast** | Minimum 4.5:1 for normal text, 3:1 for large | AA |
| **Focus Indicators** | Visible focus outlines on all elements | A |
| **Resize Text** | Content reflows up to 200% without loss of data | AA |
| **Motion Sensitivity** | Respect `prefers-reduced-motion` media query | AA |
| **Language Attributes** | `<html lang="ar-AE">` or appropriate locale | A |

## 8. Internationalization (i18n) Ready

```
Language Support:
- ar-AE (Arabic, UAE context - primary)
- en (English)
- Additional languages as needed

Translation Files Structure:
src/locales/
├─ en/
│  └─ translation.json
└─ ar-AE/
   └─ translation.json

Example translation key:
{
  "common:save": "Save",
  "common:cancel": "Cancel",
  "table:attendance:present": "Present",
  "table:attendance:absent": "Absent",
  "table:attendance:excused": "Excused"
}
```

---
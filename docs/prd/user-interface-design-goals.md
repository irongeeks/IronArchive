# User Interface Design Goals

## Overall UX Vision

IronArchive embraces a **progressive disclosure** design philosophy where complexity is revealed only when needed. The interface presents an empty state ("Add Your First Tenant") to new users, gradually introducing features like storage charts, task monitoring, and advanced settings as the system populates with data. This approach prevents user overwhelm while maintaining power-user capabilities.

The UI prioritizes **real-time feedback** to build trust: live sync progress bars with 1-second auto-refresh, instant search-as-you-type results, and immediate theme switching without page reloads. Every action provides clear, immediate visual confirmation.

Navigation follows a **logical drill-down pattern** (Dashboard → Tenant → Mailbox → Email) with consistent breadcrumb trails and back navigation, allowing users to maintain mental models of where they are in the system hierarchy.

## Key Interaction Paradigms

**Instant Feedback Loop:** All user actions trigger immediate visual responses (loading states, progress indicators, success/error notifications) to maintain engagement and confidence

**Self-Service First:** Password resets, theme changes, export jobs, and common administrative tasks are designed for user completion without MSP support escalation

**Adaptive Complexity:** Interface elements appear/disappear based on context (empty states for new installations, advanced filters only when needed, role-appropriate menu options)

**Unified Multi-Tenant View:** MSP Admins see all tenants in a single dashboard rather than requiring separate logins per client—a key differentiator from legacy solutions

**Search-Centric Workflow:** Search is prominently featured as the primary access path to archived emails, with side-panel preview enabling quick browsing without full page navigation

## Core Screens and Views

**Setup Wizard Screen:** First-run experience for creating initial MSP Admin account and configuring basic system settings

**Dashboard (Adaptive):** Primary landing page showing tenant overview, storage usage charts (when data exists), task monitoring widget, and quick actions ("Add Tenant", "Manual Backup")

**Tenant Onboarding Wizard:** Multi-step flow for Azure AD app creation, mailbox discovery, mailbox selection (with badges for Shared/Licensed/No Mailbox), and backup initiation

**Tenant Detail View:** Drill-down showing single tenant's mailboxes, storage breakdown, sync status, retention policies, and tenant-specific settings

**Mailbox Detail View:** Individual mailbox view showing email list, sync history, storage usage, and mailbox-specific settings (retention, legal hold)

**Search Interface:** Full-screen search experience with filters sidebar, result list with highlighted terms, side-panel detail view, and multi-select export actions

**Task Monitoring Dashboard:** Dedicated view for all backup jobs with filters (status, tenant, date range), expandable task details, and retry capabilities

**Global Settings (MSP Admin):** Whitelabeling configuration, user/admin management, default retention policies, notification channel setup

**Profile Settings (All Users):** Theme selection with live preview, password change, MFA setup, display name, email preferences

**Export History:** List of past export jobs with download links (unexpired), status tracking, and re-export options

## Accessibility

**Target: WCAG 2.1 AA Compliance**

- Semantic HTML with proper heading hierarchy (h1, h2, h3)
- Keyboard navigation support for all interactive elements (tab order, focus indicators)
- Screen reader compatibility with ARIA labels and live regions for dynamic content
- Sufficient color contrast ratios (4.5:1 for normal text, 3:1 for large text)
- Text resize support up to 200% without loss of functionality
- Alt text for all informational images and icons
- Form labels and error messages clearly associated with inputs

## Branding

**Whitelabeling Capabilities (MSP Admin):**

- **Custom Logo Upload:** Replace IronArchive branding with MSP logo (supports PNG, SVG, JPEG; max 2MB)
- **Custom Favicon:** Upload custom favicon for browser tabs (16x16, 32x32, 48x48 sizes)
- **Custom CSS Override:** Provide sanitized CSS injection for color scheme customization while maintaining layout integrity

**User Theme System (All Users):**

Pre-installed themes allow individual users to personalize their experience without affecting MSP branding:
- Catppuccin Mocha (dark, warm tones)
- Catppuccin Latte (light, warm tones)
- Nord (cool, minimalist)
- Cyberpunk (high contrast, neon accents)
- Dracula (dark, purple accents)
- Tokyo Night (dark, blue/purple palette)

Themes apply instantly via CSS variable switching without page refresh, providing immediate visual feedback.

## Target Device and Platforms

**Web Responsive:** Primary target platform with mobile-first responsive design

**Supported Devices:**
- Desktop/Laptop (1920x1080 and higher, optimized for 2560x1440 and 4K)
- Tablet (768px-1024px, portrait and landscape)
- Mobile (320px-767px, optimized for iPhone 12/13/14, Android flagships)

**Browser Support:**
- Desktop: Chrome/Edge 120+, Firefox 120+, Safari 17+ (last 2 major versions)
- Mobile: iOS Safari 16+, Android Chrome 120+

**Not Supported:**
- Internet Explorer (all versions)
- Browsers older than 2 years
- Embedded browsers in native apps (may work but not officially supported)

**Performance Targets:**
- Dashboard load < 2s on broadband (25 Mbps+)
- Search results < 200ms (instant-as-you-type)
- Smooth scrolling (60 FPS) on modern devices
- Graceful degradation on slower connections (show loading states, reduce auto-refresh frequency)

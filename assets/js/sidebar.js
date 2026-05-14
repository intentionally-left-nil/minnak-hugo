export default function sidebar() {
    const toggleButton = document.getElementById('toggle-button');
    const leftSidebar = document.getElementById('left-sidebar');

    if (!leftSidebar) return;

    // -------------------------------------------------------
    // Mobile sidebar toggle
    // -------------------------------------------------------
    if (toggleButton) {
        toggleButton.addEventListener('mousedown', function (e) {
            document.body.classList.toggle('open-sidebar');
            e.stopPropagation();
        });

        // Close when clicking outside the sidebar
        document.addEventListener('click', function (e) {
            if (
                document.body.classList.contains('open-sidebar') &&
                !leftSidebar.contains(e.target) &&
                e.target !== toggleButton &&
                !toggleButton.contains(e.target)
            ) {
                document.body.classList.remove('open-sidebar');
            }
        });

        // Close when focus leaves the last focusable element in sidebar
        const focusable = leftSidebar.querySelectorAll(
            'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
        );
        const lastFocusable = focusable[focusable.length - 1];
        if (lastFocusable) {
            lastFocusable.addEventListener('blur', function () {
                document.body.classList.remove('open-sidebar');
            });
        }
    }

    // -------------------------------------------------------
    // Tab switching
    // -------------------------------------------------------
    const navLinks = leftSidebar.querySelectorAll('.vertical-menu-item .nav-link');

    navLinks.forEach(function (link) {
        let mousedown = false;

        link.addEventListener('mousedown', function () {
            mousedown = true;
            const item = this.parentElement;
            if (!item.classList.contains('active-menu')) {
                leftSidebar.querySelectorAll('.vertical-menu-item.active-menu')
                    .forEach(el => el.classList.remove('active-menu'));
                item.classList.add('active-menu');
            }
        });

        // Keyboard / tab focus fallback
        link.addEventListener('focusin', function () {
            if (!mousedown) {
                const item = this.parentElement;
                if (!item.classList.contains('active-menu')) {
                    leftSidebar.querySelectorAll('.vertical-menu-item.active-menu')
                        .forEach(el => el.classList.remove('active-menu'));
                    item.classList.add('active-menu');
                }
            }
            mousedown = false;
        });
    });
}

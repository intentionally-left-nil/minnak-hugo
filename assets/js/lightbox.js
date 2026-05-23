// lightbox.js — Gallery image lightbox
//
// Intercepts clicks on .gallery .gallery-item links, preventing navigation
// to the raw image file. Opens a full-screen overlay instead with the
// 1200px preview (data-lightbox-src), keeping the original href as a
// no-JS fallback.
//
// Keyboard: Escape closes. Focus moves to the close button on open and
// returns to the triggering link on close.

export default function lightbox() {
    let lb = null;
    let triggerEl = null;

    // ── Build the overlay element (once, lazily) ────────────────────────────
    function getLightbox() {
        if (lb) return lb;

        lb = document.createElement('div');
        lb.id = 'lightbox';
        lb.className = 'lightbox';
        lb.setAttribute('role', 'dialog');
        lb.setAttribute('aria-modal', 'true');
        lb.setAttribute('aria-label', 'Image viewer');
        lb.hidden = true;
        lb.innerHTML =
            '<button class="lightbox-close" aria-label="Close image viewer">\u00d7</button>' +
            '<figure class="lightbox-figure">' +
                '<img class="lightbox-img" src="" alt="">' +
                '<figcaption class="lightbox-caption"></figcaption>' +
            '</figure>';

        document.body.appendChild(lb);

        // Click on the backdrop (not on the figure) closes the lightbox.
        lb.addEventListener('click', function (e) {
            if (e.target === lb) closeLightbox();
        });

        lb.querySelector('.lightbox-close').addEventListener('click', closeLightbox);

        return lb;
    }

    // ── Open ────────────────────────────────────────────────────────────────
    function openLightbox(src, alt, caption, trigger) {
        const overlay = getLightbox();
        const img = overlay.querySelector('.lightbox-img');
        const cap = overlay.querySelector('.lightbox-caption');

        img.src = src;
        img.alt = alt;
        cap.textContent = caption || '';
        cap.hidden = !caption;

        triggerEl = trigger || null;
        overlay.hidden = false;
        document.body.style.overflow = 'hidden';
        overlay.querySelector('.lightbox-close').focus();
    }

    // ── Close ───────────────────────────────────────────────────────────────
    function closeLightbox() {
        if (!lb || lb.hidden) return;
        lb.hidden = true;
        document.body.style.overflow = '';
        if (triggerEl) {
            triggerEl.focus();
            triggerEl = null;
        }
    }

    // ── Keyboard ─────────────────────────────────────────────────────────────
    document.addEventListener('keydown', function (e) {
        if (e.key === 'Escape') closeLightbox();
    });

    // ── Click delegation ─────────────────────────────────────────────────────
    document.addEventListener('click', function (e) {
        const link = e.target.closest('.gallery .gallery-item a');
        if (!link) return;
        e.preventDefault();

        const img = link.querySelector('img');
        const alt = img ? img.alt : '';
        const figcaption = link.closest('.gallery-item').querySelector('figcaption');
        const caption = figcaption ? figcaption.textContent.trim() : '';
        const src = link.dataset.lightboxSrc || link.href;

        openLightbox(src, alt, caption, link);
    });
}

// MiNNaK Hugo Theme — search modal behaviour
//
// <pagefind-modal> has no built-in close-on-navigate handler.  Same-page
// anchor links (e.g. heading sub-results) update the URL hash without a full
// page reload, so the modal would otherwise stay open after navigation.
// Cross-page links trigger a full reload, so the new page starts with the
// modal closed — no extra handling needed there.

// Workaround for https://github.com/pagefind/pagefind/issues/1125
export default function fixSamePageSearch() {
    const modal = document.querySelector('pagefind-modal');
    if (!modal) return;

    window.addEventListener('hashchange', () => modal.close());
}

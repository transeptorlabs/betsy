const setActiveLink = (activeLinkId) => {
    document.querySelectorAll('.nav-link').forEach(link => {
      link.classList.remove('active');
    });
    document.getElementById(activeLinkId).classList.add('active');
}

document.addEventListener("DOMContentLoaded", function() {
    // Load Accounts content on page load
    htmx.ajax('GET', '/accounts', { target: '#page-content' });

    // Set the initial active link
    setActiveLink('accounts-link');

    // Update the active link on click
    document.querySelectorAll('.nav-link').forEach(link => {
      link.addEventListener('click', function() {
        setActiveLink(this.id);
      });
    });
});
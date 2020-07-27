<template functional>
  <b-card no-body class="feedback-card">
    <header>
      <h5 class="text-capitalize">
        <p class="m-0 feedback-title text-truncate">{{props.feedback.title}}</p>
        <p class="text-muted m-0 feedback-date">{{props.feedback.update_at | date}}</p>
      </h5>
      <span
        class="text-muted text-truncate"
      >From {{props.feedback.display_name}} | {{props.feedback.creator}}</span>
    </header>
    <div class="feedback-body secondary-font scrollBar">{{props.feedback.body}}</div>
  </b-card>
</template>
<style lang="scss">
.feedback-card {
  box-shadow: 0 1px 15px rgba(0, 0, 0, 0.09);
  border-radius: 5px;
  border: none;

  header {
    margin: 0;
    display: flex;
    flex-direction: column;
    padding: 0.7rem 0.5rem 0.3rem;
    align-items: flex-start;
    background: #f8f9fa;
    border: 1px solid #dee2e6;
    border-radius: inherit;
    border-bottom-left-radius: 0;
    border-bottom-right-radius: 0;

    h5 {
      display: flex;
      width: 100%;
      justify-content: space-between;
      align-items: flex-start;
      margin-bottom: 0.5rem;

      .feedback-title {
        font-size: 1.25rem;
      }

      .feedback-date {
        width: fit-content;
        white-space: nowrap;
        padding-left: 0.5rem;
      }
    }
    span {
      font-size: 0.9rem;
    }
  }
  .feedback-body {
    padding: 0.75rem;
    max-height: 150px;
    overflow: auto;
    min-height: 80px;
    font-size: 1.1rem;
  }
}
</style>
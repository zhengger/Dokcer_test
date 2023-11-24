package main

func TestRespondsWithLove(t *testing.T) {

	pool, err := dockertest.NewPool("")
	require.NoError(t, err, "could not connect to Docker")
	
	resource, err := pool.Run("docker-gs-ping", "latest", []string{})
	require.NoError(t, err, "could not start container")
	
	t.Cleanup(func() {
		require.NoError(t, pool.Purge(resource), "failed to remove container")
	})
	
	var resp *http.Response

	err = pool.Retry(func() error {
		resp, err = http.Get(fmt.Sprint("http://localhost:", resource.GetPort("8080/tcp"), "/"))
		if err != nil {
			t.Log("container not ready, waiting...")
			return err
		}
		return nil
	})
	require.NoError(t, err, "HTTP error")
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, "HTTP status code")

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "failed to read HTTP body")

	// Finally, test the business requirement!
	require.Contains(t, string(body), "<3", "does not respond with love?")	
}
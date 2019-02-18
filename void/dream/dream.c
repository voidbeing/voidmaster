#include <linux/kernel.h>
#include <linux/module.h>
#include <linux/sched.h>
#include <asm/uaccess.h>
#include <linux/fs.h>
#include <linux/input.h>


static ssize_t dream_read(struct file *filp, char *buffer, size_t length, loff_t *offset) {
	printk("read\n");
	return length;
}

static ssize_t dream_write(struct file *filp, const char *buff, size_t length, loff_t *offset) {
	printk("write\n");
	return length;
}

static int dream_open(struct inode *inode, struct file *file) {
	printk("open\n");
	return 0;
}
static int dream_release(struct inode *inode, struct file *file) {
	printk("release\n");
	struct input_handler *handler;
	list_for_each_entry(handler, &input_handler_list, node)
		printk("%s\n", handler->name);
	return 0;
}

static struct file_operations fops = {
	.open = dream_open,
	.release = dream_release,
	.read= dream_read,
	.write= dream_write,
};

#define DEVICE_NAME "dream"
int major_version;
int init_module() {
	major_version = register_chrdev(0, DEVICE_NAME, &fops);
	if (major_version < 0) {
		printk(KERN_ALERT "Register failed, error %d.\n", major_version);
		return major_version;
	}
	printk(KERN_INFO "'mknod /dev/%s c %d 0'.\n", DEVICE_NAME, major_version);
	return 0;
}

void cleanup_module(void) {
	unregister_chrdev(major_version, DEVICE_NAME);
	printk(KERN_INFO "Bye World!\n");
}
